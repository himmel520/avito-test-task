package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

// GetTenders возвращает список тендеров
func (s *Service) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	return s.repo.GetTenders(ctx, serviceType, limit, offset)
}

// CreateTender создает новый тендер
func (s *Service) CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	return s.repo.CreateTender(ctx, tender)
}

// GetUserTenders возвращает тендеры пользователя
func (s *Service) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	return s.repo.GetUserTenders(ctx, username, limit, offset)
}

// GetTenderStatus возвращает статус тендера
func (s *Service) GetTenderStatus(ctx context.Context, tenderID string) (*models.TenderStatus, error) {
	return s.repo.GetTenderStatus(ctx, tenderID)
}

// UpdateTenderStatus обновляет статус тендера
func (s *Service) UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error) {
	return s.repo.UpdateTenderStatus(ctx, tenderID, status)
}

// EditTender редактирует существующий тендер
func (s *Service) EditTender(ctx context.Context, tenderID string, username string, tender models.TenderEdit) (*models.TenderResponse, error) {
	return s.repo.UpdateTender(ctx, tender)
}

// RollbackTender откатывает тендер к указанной версии
func (s *Service) RollbackTender(ctx context.Context, tenderID string, version int32, username string) (*models.TenderResponse, error) {
	return s.repo.RollbackTender(ctx, tenderID, version)
}
