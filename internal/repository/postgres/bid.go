package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateBid создает новое предложение и возвращает его сгенерированным ID
func (p *Postgres) CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error) {
	bidResp := &models.BidResponse{}

	err := p.DB.QueryRow(ctx, `
	insert into bid 
		(name, description, tender_id, author_type, author_id)
	values 
    	($1, $2, $3, $4, $5) 
	returning *;`, bid.Name, bid.Description, bid.TenderID, bid.AuthorType, bid.AuthorId).Scan(
		&bidResp.ID, &bidResp.Name, &bidResp.Description, &bidResp.Status, &bidResp.TenderID,
		&bidResp.AuthorType, &bidResp.AuthorID, &bidResp.Version, &bidResp.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case repository.FKViolation:
			return nil, repository.ErrBidDependencyNotFound
		case repository.UniqueConstraint:
			return nil, repository.ErrBidUnique
		}
	}

	return bidResp, err
}

// GetUserBids получает список предложений пользователя с указанным ID с пагинацией
func (p *Postgres) GetUserBids(ctx context.Context, userID string, limit, offset int32) ([]*models.BidResponse, error) {
	rows, err := p.DB.Query(ctx, `
	SELECT *
	FROM bid b 
		WHERE author_id = $1
	ORDER BY name ASC
	LIMIT $2 OFFSET $3;`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bids := []*models.BidResponse{}
	for rows.Next() {
		bid := &models.BidResponse{}
		if err := rows.Scan(
			&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
			&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
			return nil, err
		}

		bids = append(bids, bid)
	}

	return bids, err
}

// GetBidsForTender получает список предложений для указанного тендера с пагинацией и исключает предложения со статусом 'Created'
func (p *Postgres) GetBidsForTender(ctx context.Context, tenderID string, limit, offset int32) ([]*models.BidResponse, error) {
	rows, err := p.DB.Query(ctx, `
	SELECT *
	FROM bid
		WHERE tender_id = $1
		AND status != 'Created' 
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3;
	`, tenderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bids := []*models.BidResponse{}
	for rows.Next() {
		bid := &models.BidResponse{}
		if err := rows.Scan(
			&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
			&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt); err != nil {
			return nil, err
		}

		bids = append(bids, bid)
	}

	if len(bids) == 0 {
		return nil, repository.ErrBidTenderNotFound
	}

	return bids, err
}

// GetBidByID получает предложение по его ID
func (p *Postgres) GetBidByID(ctx context.Context, bidID string) (*models.BidResponse, error) {
	bid := &models.BidResponse{}
	err := p.DB.QueryRow(ctx, `
        SELECT *
        FROM bid b
        WHERE b.id = $1`, bidID).Scan(
		&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
		&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrBidNotFound
	}

	return bid, err
}

// UpdateBidStatus обновляет статус предложения по его ID
func (p *Postgres) UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error) {
	bid := &models.BidResponse{}
	err := p.DB.QueryRow(ctx, `
    UPDATE bid
		SET status = $2::bid_status
		WHERE id = $1
	returning *`, bidID, status).Scan(
		&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
		&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrBidNotFound
	}

	return bid, err
}

// EditBid обновляет предложение по его ID и сохраняет старую версию в историю
func (p *Postgres) EditBid(ctx context.Context, bidID string, bidEdit *models.BidEdit) (*models.BidResponse, error) {
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
	INSERT INTO bid_version 
		(bid_id, name, description, status, tender_id, author_type, author_id, version, created_at) 
	SELECT
		id, name, description, status, tender_id, author_type, author_id, version, created_at
	FROM bid
	WHERE id = $1;`, bidID)
	if pgCmd.RowsAffected() == 0 {
		return nil, repository.ErrBidNotFound
	}

	// Обновление предложения в основной таблице и его возврат с обновленной версией
	var keys []string
	var values []interface{}

	if bidEdit.Name != nil {
		keys = append(keys, "name=$1")
		values = append(values, bidEdit.Name)
	}

	if bidEdit.Description != nil {
		keys = append(keys, fmt.Sprintf("description=$%d", len(values)+1))
		values = append(values, bidEdit.Description)
	}

	values = append(values, bidID)
	query := fmt.Sprintf(`update bid set %s, version = version + 1 where id = $%v returning *;`, strings.Join(keys, ", "), len(values))

	bid := &models.BidResponse{}
	err = tx.QueryRow(ctx, query, values...).Scan(
		&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
		&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrBidNotFound
	}

	return bid, err
}

// RollbackBid откатывает предложение к указанной версии
func (p *Postgres) RollbackBid(ctx context.Context, bidID string, version int32) (*models.BidResponse, error) {
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return nil, err
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
	INSERT INTO bid_version 
		(bid_id, name, description, status, tender_id, author_type, author_id, version, created_at) 
	SELECT
		id, name, description, status, tender_id, author_type, author_id, version, created_at
	FROM bid
	WHERE id = $1;`, bidID)
	if pgCmd.RowsAffected() == 0 {
		return nil, repository.ErrBidNotFound
	}

	// Вытащить из истории нужную версию и обновить данные в основной таблице с инкрементом версии
	bid := &models.BidResponse{}
	err = tx.QueryRow(ctx, `
	WITH bv AS (
		SELECT
			name, description, status, tender_id, author_type, author_id, version, created_at
		FROM bid_version
		WHERE bid_id = $1 AND version = $2
	)
	UPDATE bid b
	SET
		name = bv.name,
		description = bv.description,
		status = bv.status,
		tender_id = bv.tender_id,
		author_type = bv.author_type,
		author_id = bv.author_id,
		version = b.version + 1,
		created_at = bv.created_at
	FROM bv
		WHERE b.id = $1
	RETURNING b.*;`, bidID, version).Scan(
		&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
		&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrBidORVersionNotFound
	}

	return bid, err
}
