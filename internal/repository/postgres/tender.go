package postgres

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

func (p *Postgres) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int) ([]models.Tender, error) {
	// Реализация запроса к базе данных
	return nil, nil
}

func (p *Postgres) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	// Реализация создания тендера
	return models.Tender{}, nil
}

func (p *Postgres) GetUserTenders(ctx context.Context, username string, limit, offset int) ([]models.Tender, error) {
	// Реализация получения тендеров пользователя
	return nil, nil
}

func (p *Postgres) GetTenderByID(ctx context.Context, tenderID string) (models.Tender, error) {
	// Реализация получения тендера по ID
	return models.Tender{}, nil
}

func (p *Postgres) UpdateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	// Реализация обновления тендера
	return models.Tender{}, nil
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderID string, version int) (models.Tender, error) {
	// Реализация отката тендера
	return models.Tender{}, nil
}
