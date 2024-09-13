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
	// CreateBid создает новое предложение
	CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error)
	// GetUserBids возвращает предложения пользователя по имени
	GetUserBids(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error)
	// GetBidsForTender возвращает предложения для указанного тендера
	GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error)
	// GetBidStatus возвращает статус предложения с учетом прав пользователя
	GetBidStatus(ctx context.Context, bidID string, username string) (*models.BidStatus, error)
	// UpdateBidStatus обновляет статус предложения, если пользователь является его создателем
	UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error)
	// EditBid редактирует предложение, если пользователь является его создателем
	EditBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error)
	// SubmitBidDecision подает решение по предложению и обновляет его статус при необходимости
	SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error)
	// SubmitBidFeedback добавляет отзыв к предложению, если пользователь имеет права
	SubmitBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error)
	// RollbackBid откатывает предложение к предыдущей версии, если пользователь является его создателем
	RollbackBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error)
	// GetBidReviews возвращает отзывы по предложению для тендера
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
