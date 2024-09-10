package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

func (s *Service) CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error) {
	return s.repo.CreateBid(ctx, bid)
}

func (s *Service) GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error) {
	return s.repo.GetUserBids(ctx, username, limit, offset)
}

func (s *Service) GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error) {
	return s.repo.GetBidsForTender(ctx, tenderID, username, limit, offset)
}

func (s *Service) GetBidStatus(ctx context.Context, bidID string, username string) (*models.BidStatus, error) {
	return s.repo.GetBidStatus(ctx, bidID, username)
}

func (s *Service) UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error) {
	return s.repo.UpdateBidStatus(ctx, bidID, username, status)
}

func (s *Service) EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error) {
	return s.repo.EditBid(ctx, bidID, username, bid)
}

func (s *Service) SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error) {
	return s.repo.SubmitBidDecision(ctx, bidID, username, decision)
}

func (s *Service) SubmitBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error) {
	return s.repo.SubmitBidFeedback(ctx, bidID, username, feedback)
}

func (s *Service) RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error) {
	return s.repo.RollbackBid(ctx, bidID, username, version)
}

func (s *Service) GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReview, error) {
	return s.repo.GetBidReviews(ctx, tenderID, authorUsername, requesterUsername, limit, offset)
}
