package service

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

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

func (s *Service) GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error) {
	// проверяет существование пользователя
	userId, err := s.repo.GetUserIDByName(ctx, username)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserBids(ctx, userId, limit, offset)
}

func (s *Service) GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error) {
	// проверяет, является ли пользователь ответственным за организацию тендера
	if err := s.repo.IsUserResponsibleForTender(ctx, tenderID, username); err != nil {
		return nil, err
	}

	return s.repo.GetBidsForTender(ctx, tenderID, limit, offset)
}

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

func (s *Service) UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error) {
	// проверяет, является ли пользователь создателем предложения
	if err := s.repo.IsBidCreatorByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.UpdateBidStatus(ctx, bidID, username, status)
}

func (s *Service) EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error) {
	// проверяет, является ли пользователь создателем предложения
	if err := s.repo.IsBidCreatorByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.EditBid(ctx, bidID, bid)
}

func (s *Service) SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error) {
	return s.repo.SubmitBidDecision(ctx, bidID, username, decision)
}

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

func (s *Service) RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error) {
	// проверяет, является ли пользователь создателем предложения
	if err := s.repo.IsBidCreatorByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.RollbackBid(ctx, bidID, version)
}

func (s *Service) GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReviewResponse, error) {
	// проверяет, относится ли пользователь к организации тендера
	if err := s.repo.IsUserResponsibleForTender(ctx, tenderID, requesterUsername); err != nil {
		// пользователь не относится к создателем тендера
		return nil, err
	}

	return s.repo.GetBidReviews(ctx, tenderID, authorUsername, limit, offset)
}
