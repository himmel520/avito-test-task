package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

// GetTenders возвращает список тендеров
func (s *Service) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int) ([]models.Tender, error) {
	return s.repo.GetTenders(ctx, serviceType, limit, offset)
}

// CreateTender создает новый тендер
func (s *Service) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	return s.repo.CreateTender(ctx, tender)
}

// GetUserTenders возвращает тендеры пользователя
func (s *Service) GetUserTenders(ctx context.Context, username string, limit, offset int) ([]models.Tender, error) {
	return s.repo.GetUserTenders(ctx, username, limit, offset)
}

// GetTenderStatus возвращает статус тендера
func (s *Service) GetTenderStatus(ctx context.Context, tenderID string) (models.TenderStatus, error) {
	tender, err := s.repo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return "", err
	}
	return tender.Status, nil
}

// UpdateTenderStatus обновляет статус тендера
func (s *Service) UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (models.Tender, error) {
	tender, err := s.repo.GetTenderByID(ctx, tenderID)
	if err != nil {
		return models.Tender{}, err
	}
	tender.Status = status
	return s.repo.UpdateTender(ctx, tender)
}

// EditTender редактирует существующий тендер
func (s *Service) EditTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	return s.repo.UpdateTender(ctx, tender)
}

// RollbackTender откатывает тендер к указанной версии
func (s *Service) RollbackTender(ctx context.Context, tenderID string, version int) (models.Tender, error) {
	return s.repo.RollbackTender(ctx, tenderID, version)
}
