package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	// формируем фильтр по enum service_type
	var filter string
	if len(serviceType) != 0 {
		var types []string
		for _, stype := range serviceType {
			types = append(types, fmt.Sprintf("'%v'::service_type", stype))
		}

		filter = fmt.Sprintf(
			"AND service_type = ANY (ARRAY[%s])",
			strings.Join(types, ","),
		)
	}

	query := fmt.Sprintf(`
	SELECT 
		id, name, description, service_type, status, organization_id, version, created_at
	FROM tender
	WHERE status = 'Published'
	%s
	ORDER BY name ASC 
	LIMIT $1 OFFSET $2;`, filter)

	rows, err := p.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := []*models.TenderResponse{}
	for rows.Next() {
		tender := &models.TenderResponse{}
		if err := rows.Scan(
			&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
			&tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (p *Postgres) CreateTender(ctx context.Context, tender *models.TenderCreate, employeeId string) (*models.TenderResponse, error) {
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	tenderResp := &models.TenderResponse{}
	err = tx.QueryRow(ctx, `
	insert into tender 
		(name, description, service_type, organization_id) 
	values ($1, $2, $3, $4) returning *;`,
		tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID).Scan(
		&tenderResp.ID, &tenderResp.Name, &tenderResp.Description,
		&tenderResp.ServiceType, &tenderResp.Status, &tenderResp.OrganizationID, &tenderResp.Version, &tenderResp.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == repository.FKViolation {
				return nil, repository.ErrOrganizationDepencyNotFound
			}
		}
		return nil, fmt.Errorf("failed to insert tender: %w", err)
	}

	_, err = tx.Exec(ctx, `
	insert into tender_creator 
		(creator_id, tender_id) 
	values ($1, $2)`, employeeId, tenderResp.ID)

	return tenderResp, err
}

func (p *Postgres) GetUserTenders(ctx context.Context, userId string, limit, offset int32) ([]*models.TenderResponse, error) {
	rows, err := p.DB.Query(ctx, `
	SELECT 
		t.*
	FROM tender t
	JOIN tender_creator tc ON t.id = tc.tender_id
	WHERE 
		tc.creator_id = $3  
	ORDER BY t.name ASC  
	LIMIT $1  
	OFFSET $2; `, limit, offset, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := []*models.TenderResponse{}
	for rows.Next() {
		tender := &models.TenderResponse{}
		if err := rows.Scan(
			&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
			&tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (p *Postgres) GetTenderStatus(ctx context.Context, tenderID string) (*models.TenderStatus, *models.OrganizationID, error) {
	var status *models.TenderStatus
	var organizationID *models.OrganizationID

	err := p.DB.QueryRow(ctx, `
	select 
		status, organization_id
	from tender 
	where id = $1;`, tenderID).Scan(&status, &organizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, repository.ErrTenderNotFound
	}

	return status, organizationID, err
}

func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error) {
	tender := &models.TenderResponse{}

	err := p.DB.QueryRow(ctx, `
	UPDATE tender
	SET status = $2::tender_status
	WHERE id = $1
	returning *;`, tenderID, status).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderNotFound
	}

	return tender, err
}

func (p *Postgres) UpdateTender(ctx context.Context, tenderID string, tenderEdit *models.TenderEdit) (*models.TenderResponse, error) {
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Добавляем текущую версию в историю
	pgCmd, err := tx.Exec(ctx, `
	INSERT INTO tender_version 
		(tender_id, name, description, service_type, status, organization_id, version, created_at) 
	SELECT
    	id, name, description, service_type, status, organization_id, version, created_at
	FROM tender
	WHERE id = $1;`, tenderID)
	if pgCmd.RowsAffected() == 0 {
		return nil, repository.ErrTenderNotFound
	}
	
	// Обновление тендера в основной таблице и его возврат с обновленной версией
	var keys []string
	var values []interface{}

	if tenderEdit.Name != nil {
		keys = append(keys, "name=$1")
		values = append(values, tenderEdit.Name)
	}

	if tenderEdit.Description != nil {
		keys = append(keys, fmt.Sprintf("description=$%d", len(values)+1))
		values = append(values, tenderEdit.Description)
	}

	if tenderEdit.ServiceType != nil {
		keys = append(keys, fmt.Sprintf("service_type=$%d::service_type", len(values)+1))
		values = append(values, tenderEdit.ServiceType)
	}

	values = append(values, tenderID)
	query := fmt.Sprintf(`update tender set %s, version = version + 1 where id = $%v returning *;`, strings.Join(keys, ", "), len(values))

	tender := &models.TenderResponse{}
	err = tx.QueryRow(ctx, query, values...).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderNotFound
	}

	return tender, err
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error) {
	// Реализация отката тендера
	return nil, nil
}

func (p *Postgres) СheckOrganizationPermission(ctx context.Context, organizationID *models.OrganizationID, username string) (string, error) {
	var existsRelation bool
	var creatorId string

	err := p.DB.QueryRow(ctx, `
	SELECT
		e.id AS user_id,
		EXISTS (
			SELECT 1
			FROM organization_responsible
			WHERE user_id = e.id
			AND organization_id = $2
		) AS exists_relation
	FROM employee e
	WHERE username = $1;`, username, organizationID).Scan(&creatorId, &existsRelation)

	switch {
	// пользователь не существует или некорректен.
	case errors.Is(err, pgx.ErrNoRows):
		return "", repository.ErrUserNotExist
	// нет связи пользователь и организация
	case !existsRelation:
		return "", repository.ErrRelationNotExist
	}

	return creatorId, err
}

func (p *Postgres) IsTenderCreator(ctx context.Context, tenderId, username string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM tender_creator tc
		WHERE tc.creator_id = e.id AND tc.tender_id = $2
	) AS is_creator
	FROM employee e
	WHERE e.username = $1;`, username, tenderId).Scan(&isCreator)

	switch {
	// пользователь не существует или некорректен.
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	// нет связи пользователь и тендер
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) GetUserIDByName(ctx context.Context, username string) (string, error) {
	var userId string
	err := p.DB.QueryRow(ctx, `
	select 
		id 
	from employee e 
	where e.username = $1;
	`, username).Scan(&userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", repository.ErrUserNotExist
	}

	return userId, err
}
