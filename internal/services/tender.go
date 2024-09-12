package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

// GetTenders возвращает список тендеров с применением фильтрации и пагинации
func (s *Service) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	return s.repo.GetTenders(ctx, serviceType, limit, offset)
}

// CreateTender создает новый тендер после проверки прав пользователя
func (s *Service) CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	// проверка права пользователя на создание тендера
	employeeId, err := s.repo.СheckOrganizationPermission(ctx, &tender.OrganizationID, tender.CreatorUsername)
	if err != nil {
		// пользователь не имеет доступа
		return nil, err
	}

	return s.repo.CreateTender(ctx, tender, employeeId)
}

// GetUserTenders возвращает тендеры, созданные указанным пользователем
func (s *Service) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	// проверяет существование пользователя
	userId, err := s.repo.GetUserIDByName(ctx, username)
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserTenders(ctx, userId, limit, offset)
}

// GetTenderStatus возвращает статус тендера после проверки прав пользователя
func (s *Service) GetTenderStatus(ctx context.Context, tenderID, username string) (*models.TenderStatus, error) {
	// получить статус и id организации
	status, organizationID, err := s.repo.GetTenderStatus(ctx, tenderID)
	if err != nil {
		return nil, err
	}

	// проверяет права пользователя на доступ к информации о тендере
	if _, err := s.repo.СheckOrganizationPermission(ctx, organizationID, username); err != nil {
		return nil, err
	}

	return status, nil
}

// UpdateTenderStatus обновляет статус тендера после проверки прав пользователя
func (s *Service) UpdateTenderStatus(ctx context.Context, tenderID, username string, status models.TenderStatus) (*models.TenderResponse, error) {
	// проверяет, является ли пользователь создателем тендер
	if err := s.repo.IsTenderCreator(ctx, tenderID, username); err != nil {
		return nil, err
	}
	
	return s.repo.UpdateTenderStatus(ctx, tenderID, status)
}

// EditTender редактирует существующий тендер после проверки прав пользователя
func (s *Service) EditTender(ctx context.Context, tenderID string, username string, tender *models.TenderEdit) (*models.TenderResponse, error) {
	// проверяет, является ли пользователь создателем тендера
	if err := s.repo.IsTenderCreator(ctx, tenderID, username); err != nil {
		return nil, err
	}

	return s.repo.UpdateTender(ctx, tenderID ,tender)
}

// RollbackTender откатывает тендер к указанной версии после проверки прав пользователя
func (s *Service) RollbackTender(ctx context.Context, tenderID string, version int32, username string) (*models.TenderResponse, error) {
	// проверяет, является ли пользователь создателем тендера
	if err := s.repo.IsTenderCreator(ctx, tenderID, username); err != nil {
		return nil, err
	}
	
	return s.repo.RollbackTender(ctx, tenderID, version)
}
