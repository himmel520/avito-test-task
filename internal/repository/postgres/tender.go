package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

// GetTenders возвращает список тендеров с заданными фильтрами и пагинацией
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

// CreateTender создает новый тендер
func (p *Postgres) CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	tenderResp := &models.TenderResponse{}

	err := p.DB.QueryRow(ctx, `
	insert into tender 
		(name, description, service_type, organization_id, creator_username) 
	values ($1, $2, $3, $4, $5) returning *;`,
		tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername).Scan(
		&tenderResp.ID, &tenderResp.Name, &tenderResp.Description, &tenderResp.ServiceType,
		&tenderResp.Status, &tenderResp.OrganizationID, &tenderResp.Version, &tenderResp.CreatedAt, &tenderResp.CreatorUsername)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == repository.FKViolation {
				return nil, repository.ErrOrganizationDepencyNotFound
			}
		}
		return nil, fmt.Errorf("failed to insert tender: %w", err)
	}

	return tenderResp, err
}

// GetUserTenders возвращает список тендеров, созданных конкретным пользователем
func (p *Postgres) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	rows, err := p.DB.Query(ctx, `
	SELECT *
	FROM tender
		WHERE creator_username = $1
	ORDER BY name ASC
	LIMIT $2 OFFSET $3;`, username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := []*models.TenderResponse{}
	for rows.Next() {
		tender := &models.TenderResponse{}
		if err := rows.Scan(
			&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
			&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername); err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil
}

// GetTenderStatus возвращает статус и id организации для указанного тендера
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

// UpdateTenderStatus обновляет статус тендера и возвращает обновленный тендер
func (p *Postgres) UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error) {
	tender := &models.TenderResponse{}

	err := p.DB.QueryRow(ctx, `
	UPDATE tender
	SET status = $2::tender_status
	WHERE id = $1
	returning *;`, tenderID, status).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderNotFound
	}

	return tender, err
}

// UpdateTender обновляет данные тендера и добавляет запись о новой версии в историю
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
		(tender_id, name, description, service_type, status, organization_id, version, created_at, creator_username) 
	SELECT
    	id, name, description, service_type, status, organization_id, version, created_at, creator_username
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
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderNotFound
	}

	return tender, err
}

// RollbackTender откатывает тендер к указанной версии из истории и возвращает обновленный тендер
func (p *Postgres) RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error) {
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

	log.Println(tenderID)
	// Добавляем текущую версию в историю
	pgCmd, err := tx.Exec(ctx, `
	INSERT INTO tender_version 
		(tender_id, name, description, service_type, status, organization_id, version, created_at, creator_username) 
	SELECT
    	id, name, description, service_type, status, organization_id, version, created_at, creator_username
	FROM tender
	WHERE id = $1;`, tenderID)
	if pgCmd.RowsAffected() == 0 {
		return nil, repository.ErrTenderNotFound
	}

	// Вытащить из истории нужную версию и обновить данные в основной таблице с инкрементом версии
	tender := &models.TenderResponse{}
	err = tx.QueryRow(ctx, `
	with tv as (
		select
			name, description, service_type, status, organization_id, version, created_at, creator_username
		from tender_version
			where tender_id = $1 and version = $2
	)
	update tender t
	set
		name = tv.name,
		description = tv.description,
		service_type = tv.service_type,
		status = tv.status,
		organization_id = tv.organization_id,
		version = t.version + 1,
		created_at = tv.created_at,
		creator_username = tv.creator_username
	from tv
		where t.id = $1 
	returning t.*;`, tenderID, version).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderORVersionNotFound
	}

	return tender, err
}
