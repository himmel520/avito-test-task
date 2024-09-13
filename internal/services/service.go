package service

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/sirupsen/logrus"
)

type TenderService interface {
	// GetTenders возвращает список тендеров с применением фильтрации и пагинации
	GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error)
	// CreateTender создает новый тендер после проверки прав пользователя
	CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error)
	// GetUserTenders возвращает тендеры, созданные указанным пользователем
	GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error)
	// GetTenderStatus возвращает статус тендера после проверки прав пользователя
	GetTenderStatus(ctx context.Context, tenderID, username string) (*models.TenderStatus, error)
	// UpdateTenderStatus обновляет статус тендера после проверки прав пользователя
	UpdateTenderStatus(ctx context.Context, tenderID, username string, status models.TenderStatus) (*models.TenderResponse, error)
	// EditTender редактирует существующий тендер после проверки прав пользователя
	EditTender(ctx context.Context, tenderID string, username string, tender *models.TenderEdit) (*models.TenderResponse, error)
	// RollbackTender откатывает тендер к указанной версии после проверки прав пользователя
	RollbackTender(ctx context.Context, tenderID string, version int32, username string) (*models.TenderResponse, error)
}

type BidService interface {
	CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error)
	GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidStatus(ctx context.Context, bidID string, username string) (*models.BidStatus, error)
	UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error)
	EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error)
	SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error)
	SubmitBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error)
	RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error)
	GetBidReviews(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReviewResponse, error)
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
