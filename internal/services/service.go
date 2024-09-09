package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/sirupsen/logrus"
)

type TenderService interface {
	GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int) ([]models.Tender, error)
	CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error)
	GetUserTenders(ctx context.Context, username string, limit, offset int) ([]models.Tender, error)
	GetTenderStatus(ctx context.Context, tenderID string) (models.TenderStatus, error)
	UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (models.Tender, error)
	EditTender(ctx context.Context, tender models.Tender) (models.Tender, error)
	RollbackTender(ctx context.Context, tenderID string, version int) (models.Tender, error)
}

type BidService interface {
	CreateBid(ctx context.Context, bid models.Bid) (models.Bid, error)
	GetUserBids(ctx context.Context, username string, limit, offset int) ([]models.Bid, error)
	GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int) ([]models.Bid, error)
	GetBidByID(ctx context.Context, bidID string) (models.Bid, error)
	UpdateBidStatus(ctx context.Context, bidID string, status models.BidStatus) (models.Bid, error)
	EditBid(ctx context.Context, bid models.Bid) (models.Bid, error)
	SubmitBidDecision(ctx context.Context, bidID string, decision models.BidDecision) (models.Bid, error)
	SubmitBidFeedback(ctx context.Context, bidID string, feedback models.BidFeedback) (models.Bid, error)
	RollbackBid(ctx context.Context, bidID string, version int) (models.Bid, error)
	GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int) ([]models.BidReview, error)
}

type Service struct {
	repo repository.Repository
	log  *logrus.Logger
}

func New(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}
