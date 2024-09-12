package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

type TenderRepository interface {
	// GetTenders возвращает список тендеров с заданными фильтрами и пагинацией
	GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error)
	// CreateTender создает новый тендер и связывает его с сотрудником
	CreateTender(ctx context.Context, tender *models.TenderCreate, employeeId string) (*models.TenderResponse, error)
	// GetUserTenders возвращает список тендеров, созданных конкретным пользователем
	GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error)
	// GetTenderStatus возвращает статус и организацию для указанного тендера
	GetTenderStatus(ctx context.Context, tenderID string) (*models.TenderStatus, *models.OrganizationID, error)
	// UpdateTenderStatus обновляет статус тендера и возвращает обновленный тендер
	UpdateTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error)
	// UpdateTender обновляет данные тендера и добавляет запись о новой версии в историю
	UpdateTender(ctx context.Context, tenderID string, tenderEdit *models.TenderEdit) (*models.TenderResponse, error)
	// RollbackTender откатывает тендер к указанной версии из истории и возвращает обновленный тендер
	RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error)

	// СheckOrganizationPermission проверяет, имеет ли пользователь права доступа к организации
	СheckOrganizationPermission(ctx context.Context, organizationID *models.OrganizationID, username string) (string, error)
	// IsTenderCreator проверяет, является ли пользователь создателем указанного тендера
	IsTenderCreator(ctx context.Context, tenderId, username string) error
	// GetUserIDByName возвращает id пользователя по его имени
	GetUserIDByName(ctx context.Context, username string) (string, error)
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
