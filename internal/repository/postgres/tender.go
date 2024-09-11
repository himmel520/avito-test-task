package postgres

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	// Реализация запроса к базе данных
	return nil, nil
}

func (p *Postgres) CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	tenderResp := &models.TenderResponse{}

	err := p.DB.QueryRow(ctx, `
	insert into tender 
		(name, description, service_type, organization_id) 
	values ($1, $2, $3, $4) returning *;`, 
	tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID).Scan(
		&tenderResp.ID, &tenderResp.Name, &tenderResp.Description, 
		&tenderResp.ServiceType, &tenderResp.Status, &tenderResp.OrganizationID, &tenderResp.Version, &tenderResp.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == repository.FKViolation {
			return nil, repository.ErrOrganizationDepencyNotFound
		}
	}

	return tenderResp, err
}

func (p *Postgres) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	// Реализация получения тендеров пользователя
	return nil, nil
}

func (p *Postgres) GetTenderStatus(ctx context.Context, tenderID string) (*models.TenderStatus, error) {
	return nil, nil
}

func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error) {
	return nil, nil
}

func (p *Postgres) UpdateTender(ctx context.Context, tender models.TenderEdit) (*models.TenderResponse, error) {
	// Реализация обновления тендера
	return nil, nil
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error) {
	// Реализация отката тендера
	return nil, nil
}

func (p *Postgres) СheckOrganizationPermission(ctx context.Context, organizationID, username string) error {
	var existsRelation bool
	err := p.DB.QueryRow(ctx, `
	SELECT
		EXISTS (
			SELECT 1
			FROM organization_responsible
			WHERE user_id = e.id
			AND organization_id = $2
		) AS exists_relation
	FROM employee e
	WHERE username = $1;`, username, organizationID).Scan(&existsRelation)
	switch {
	// пользователь не существует или некорректен.
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и организация
	case !existsRelation:
		return repository.ErrRelationNotExist
	}

	// existsRelation = true; err = nil/error
	return err
}
