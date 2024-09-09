package postgres

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

func (p *Postgres) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	// Реализация запроса к базе данных
	return nil, nil
}

func (p *Postgres) CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	// Реализация создания тендера
	return nil, nil
}

func (p *Postgres) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	// Реализация получения тендеров пользователя
	return nil, nil
}

func (p *Postgres) GetTenderStatus(ctx context.Context, tenderID string) (*models.TenderStatus, error){
	return nil, nil
}

func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error){
	return nil, nil
}

func (p *Postgres) UpdateTender(ctx context.Context, tender models.TenderEdit) (*models.TenderResponse, error) {
	// Реализация обновления тендера
	return nil, nil
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderID string, version int) (*models.TenderResponse, error) {
	// Реализация отката тендера
	return nil, nil
}
