package postgres

import (
	"context"
	"errors"
	"log"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/jackc/pgx/v5"
)

// SubmitBidDecision записывает решение по предложению и возвращает обновленное предложение
func (p *Postgres) SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error) {
	bid := &models.BidResponse{}
	err := p.DB.QueryRow(ctx, `
	WITH inserted AS (
		INSERT INTO bid_decision (bid_id, user_id, decision)
		VALUES ($1, (SELECT id FROM employee WHERE username = $2), $3)
		RETURNING bid_id
	)
	SELECT b.*
	FROM bid b
	JOIN inserted i ON b.id = i.bid_id;`, bidID, username, decision).Scan(
		&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID,
		&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrBidNotFound
	}

	log.Println(err)

	return bid, err
}

// CountResponsibleByBid возвращает количество ответственных пользователей для указанного предложения
func (p *Postgres) CountResponsibleByBid(ctx context.Context, bidID string) (int, error) {
	var count int

	err := p.DB.QueryRow(ctx, `
	with org as (
		select t.organization_id
			from tender t
			join bid b ON t.id = b.tender_id
				where b.id = $1
	)
	select COUNT(orr.user_id) as user_count
		from organization_responsible orr
		join org o on orr.organization_id = o.organization_id;`, bidID).Scan(&count)

	return count, err
}

// CountApprovedDecisions возвращает количество одобренных решений для указанного bidID
func (p *Postgres) CountApprovedDecisions(ctx context.Context, bidID string) (int, error) {
    var count int
    err := p.DB.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM bid_decision
        WHERE bid_id = $1 AND decision = $2`, bidID, models.BidDecisionApproved).Scan(&count)
    
    return count, err
}
