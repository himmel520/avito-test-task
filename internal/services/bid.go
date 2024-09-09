package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

func (s *Service) CreateBid(ctx context.Context, bid models.Bid) (models.Bid, error) {
	return s.repo.CreateBid(ctx, bid)
}

func (s *Service) GetUserBids(ctx context.Context, username string, limit, offset int) ([]models.Bid, error) {
	return s.repo.GetUserBids(ctx, username, limit, offset)
}

func (s *Service) GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int) ([]models.Bid, error) {
	return s.repo.GetBidsForTender(ctx, tenderID, username, limit, offset)
}

func (s *Service) GetBidByID(ctx context.Context, bidID string) (models.Bid, error) {
	return s.repo.GetBidByID(ctx, bidID)
}

func (s *Service) UpdateBidStatus(ctx context.Context, bidID string, status models.BidStatus) (models.Bid, error) {
	return s.repo.UpdateBidStatus(ctx, bidID, status)
}

func (s *Service) EditBid(ctx context.Context, bid models.Bid) (models.Bid, error) {
	return s.repo.EditBid(ctx, bid)
}

func (s *Service) SubmitBidDecision(ctx context.Context, bidID string, decision models.BidDecision) (models.Bid, error) {
	return s.repo.SubmitBidDecision(ctx, bidID, decision)
}

func (s *Service) SubmitBidFeedback(ctx context.Context, bidID string, feedback models.BidFeedback) (models.Bid, error) {
	return s.repo.SubmitBidFeedback(ctx, bidID, feedback)
}

func (s *Service) RollbackBid(ctx context.Context, bidID string, version int) (models.Bid, error) {
	return s.repo.RollbackBid(ctx, bidID, version)
}

func (s *Service) GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int) ([]models.BidReview, error) {
	return s.repo.GetBidReviews(ctx, tenderID, authorUsername, requesterUsername, limit, offset)
}
