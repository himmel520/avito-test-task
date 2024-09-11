package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

type TenderRepository interface {
	GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error)
	CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error)
	GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error)
	GetTenderStatus(ctx context.Context, tenderID string) (*models.TenderStatus, error)
	UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error)
	UpdateTender(ctx context.Context, tender models.TenderEdit) (*models.TenderResponse, error)
	RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error)

	СheckOrganizationPermission(ctx context.Context, organizationID, username string) error
}

type BidRepository interface {
	CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error)
	GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidStatus(ctx context.Context, bidID string, username string) (*models.BidStatus, error)
	UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error)
	EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error)
	SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error)
	SubmitBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error)
	RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error)
	GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReview, error)
}

type Repository interface {
	TenderRepository
	BidRepository
}
