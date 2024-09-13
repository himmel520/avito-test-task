package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

// CreateBid создает новое предложение, проверяя статус тендера и права пользователя
func (s *Service) CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error) {
	// Проверка статуса тендера
	if err := s.repo.IsTenderPudlished(ctx, bid.TenderID); err != nil {
		return nil, err
	}

	// проверяет, является ли пользователь ответственным за любую организацию
	if err := s.repo.IsUserResponsible(ctx, bid.AuthorId); err != nil {
		return nil, err
	}

	return s.repo.CreateBid(ctx, bid)
}

// GetUserBids возвращает предложения пользователя по его имени
func (s *Service) GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error) {
	// проверяет существование пользователя
	userId, err := s.repo.GetUserIDByName(ctx, username)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserBids(ctx, userId, limit, offset)
}

// GetBidsForTender возвращает предложения для указанного тендера
func (s *Service) GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error) {
	// проверяет, является ли пользователь ответственным за организацию тендера
	if err := s.repo.IsUserResponsibleForTender(ctx, tenderID, username); err != nil {
		return nil, err
	}

	return s.repo.GetBidsForTender(ctx, tenderID, limit, offset)
}

// GetBidStatus возвращает статус предложения с учетом прав пользователя
func (s *Service) GetBidStatus(ctx context.Context, bidID string, username string) (*models.BidStatus, error) {
	bid, err := s.repo.GetBidByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// относится ли пользователь к организации автора bid
	err = s.repo.IsUserResponsibleForAuthorBid(ctx, bidID, username)
	switch {
	case err == nil:
		// пользователь относится к организации автора bid
		return &bid.Status, nil
	case !errors.Is(err, repository.ErrRelationNotExist):
		// произошла другая ошибка
		return nil, err
	}

	// проверяет, относится ли пользователь к организации тендера
	if err := s.repo.IsUserResponsibleForTender(ctx, bid.TenderID, username); err != nil {
		// пользователь не относится к организации автора bid и создателем тендера
		return nil, err
	}

	// у создателя тендера нет доступа к bid со статусом created
	if bid.Status == models.BidStatusCreated {
		return nil, repository.ErrRelationNotExist
	}

	return &bid.Status, nil
}

// UpdateBidStatus обновляет статус предложения, если пользователь является его создателем
func (s *Service) UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error) {
	// проверяет, является ли пользователь создателем предложения
	if err := s.repo.IsBidCreatorByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.UpdateBidStatus(ctx, bidID, username, status)
}

// EditBid редактирует предложение, если пользователь является его создателем
func (s *Service) EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error) {
	// проверяет, является ли пользователь создателем предложения
	if err := s.repo.IsBidCreatorByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.EditBid(ctx, bidID, bid)
}

// SubmitBidDecision подает решение по предложению и обновляет его статус при необходимости
func (s *Service) SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error) {
	// доступно только ответственным за организацию, связанной с тендером.
	if err := s.repo.IsUserResponsibleForTenderByBidID(ctx, bidID, username); err != nil {
		// пользователь не относится к организации тендера
		return nil, err
	}

	// создаем запись о решении
	bid, err := s.repo.SubmitBidDecision(ctx, bidID, username, decision)
	if err != nil {
		return nil, err
	}

	// если decision = "Rejected"
	if *decision == models.BidDecisionRejected {
		// сразу отклоняем преложение
		return s.repo.UpdateBidStatus(ctx, bidID, username, &models.BidStatusCanceled)
	}

	quorum, err := s.getQuorum(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// Получаем количество решений "Approved"
	approvedCount, err := s.repo.CountApprovedDecisions(ctx, bidID)
	if err != nil {
		s.log.Info(err)
		return nil, err
	}

	// Если количество решений "Approved" превысило кворум
	if approvedCount >= quorum {
		// обновляем статус bid
		bid, err := s.repo.UpdateBidStatus(ctx, bidID, username, &models.BidStatusApproved)
		if err != nil {
			s.log.Info(err)
			return nil, err
		}
		// закрываем тендер
		_, err = s.repo.UpdateTenderStatus(ctx, bid.TenderID, models.TenderStatusClosed)
		s.log.Info(err)
		return bid, err

	}

	return bid, nil
}

// SubmitBidFeedback добавляет отзыв к предложению, если пользователь имеет права
func (s *Service) SubmitBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error) {
	bid, err := s.repo.GetBidByID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	// проверяет, относится ли пользователь к организации тендера
	if err := s.repo.IsUserResponsibleForTender(ctx, bid.TenderID, username); err != nil {
		// пользователь не относится к создателем тендера
		return nil, err
	}

	return bid, s.repo.SubmitBidFeedback(ctx, bidID, feedback)
}

// RollbackBid откатывает предложение к предыдущей версии, если пользователь является его создателем
func (s *Service) RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error) {
	// проверяет, является ли пользователь создателем предложения
	if err := s.repo.IsBidCreatorByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.RollbackBid(ctx, bidID, version)
}

// GetBidReviews возвращает отзывы по предложению для тендера
func (s *Service) GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReviewResponse, error) {
	// проверяет, относится ли пользователь к организации тендера
	if err := s.repo.IsUserResponsibleForTender(ctx, tenderID, requesterUsername); err != nil {
		// пользователь не относится к создателем тендера
		return nil, err
	}

	return s.repo.GetBidReviews(ctx, tenderID, authorUsername, limit, offset)
}

// getQuorum возвращает количество голосов, необходимое для одобрения предложения
func (s *Service) getQuorum(ctx context.Context, bidID string) (int, error) {
	// Получаем количество ответственных за организацию
	count, err := s.repo.CountResponsibleByBid(ctx, bidID)
	// Рассчитываем кворум
	return min(3, count), err
}
