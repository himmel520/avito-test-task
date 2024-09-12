package postgres

import (
	"context"
	"errors"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgx/v5/pgconn"
)

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

func (p *Postgres) GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error) {
	// Реализация получения предложений пользователя
	return nil, nil
}

func (p *Postgres) GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error) {
	// Реализация получения предложений для тендера
	return nil, nil
}

func (p *Postgres) GetBidStatus(ctx context.Context, bidID string, username string) (*models.BidStatus, error) {
	return nil, nil
}

func (p *Postgres) UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error) {
	// Реализация обновления статуса предложения
	return nil, nil
}

func (p *Postgres) EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error) {
	// Реализация редактирования предложения
	return nil, nil
}

func (p *Postgres) SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error) {
	// Реализация отправки решения по предложению
	return nil, nil
}

func (p *Postgres) SubmitBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error) {
	// Реализация отправки отзыва по предложению
	return nil, nil
}

func (p *Postgres) RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error) {
	// Реализация отката версии предложения
	return nil, nil
}

func (p *Postgres) GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReview, error) {
	// Реализация получения отзывов на предложения
	return nil, nil
}
