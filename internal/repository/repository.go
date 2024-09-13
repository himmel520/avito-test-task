package repository

import (
	"context"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
)

type TenderRepository interface {
	// GetTenders возвращает список тендеров с заданными фильтрами и пагинацией
	GetTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error)
	// CreateTender создает новый тендер
	CreateTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error)
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
	СheckOrganizationPermission(ctx context.Context, organizationID *models.OrganizationID, username string) error
	// IsTenderCreatorByName проверяет, является ли пользователь создателем указанного тендера по username
	IsTenderCreatorByName(ctx context.Context, tenderId, creatorUsername string) error
	// IsUserResponsible проверяет, имеет ли пользователь права доступа к любой организации
	IsUserResponsible(ctx context.Context, userId string) error
	// GetUserIDByName возвращает id пользователя по его имени
	GetUserIDByName(ctx context.Context, username string) (string, error)
	// IsTenderPublished проверяет, опубликован ли тендер
	IsTenderPudlished(ctx context.Context, tenderID string) error
}

type BidRepository interface {
	// CreateBid создает новое предложение и возвращает его в ответе
	CreateBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error)
	// GetUserBids возвращает список предложений пользователя с указанным ID, с учетом лимита и смещения
	GetUserBids(ctx context.Context, userID string, limit, offset int32) ([]*models.BidResponse, error)
	// GetBidsForTender возвращает список предложений для указанного тендера с учетом лимита и смещения
	GetBidsForTender(ctx context.Context, tenderID string, limit, offset int32) ([]*models.BidResponse, error)
	// GetBidByID возвращает предложение по его ID
	GetBidByID(ctx context.Context, bidID string) (*models.BidResponse, error)
	// UpdateBidStatus обновляет статус предложения и возвращает обновленное предложение
	UpdateBidStatus(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error)
	// EditBid редактирует предложение и возвращает обновленное предложение
	EditBid(ctx context.Context, bidID string, bidEdit *models.BidEdit) (*models.BidResponse, error)
	// SubmitBidDecision отправляет решение по предложению и возвращает обновленное предложение
	SubmitBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error)
	// SubmitBidFeedback отправляет отзыв по предложению
	SubmitBidFeedback(ctx context.Context, bidID string, feedback *models.BidFeedback) error
	// RollbackBid откатывает предложение к указанной версии и возвращает обновленное предложение
	RollbackBid(ctx context.Context, bidID string, version int32) (*models.BidResponse, error)
	// GetBidReviews возвращает отзывы по предложениям для указанного тендера и автора, с учетом лимита и смещения
	GetBidReviews(ctx context.Context, tenderID, authorUsername string, limit, offset int32) ([]*models.BidReviewResponse, error)
	
	// IsBidCreatorByName проверяет, является ли пользователь создателем предложения по имени пользователя
	IsBidCreatorByName(ctx context.Context, bidID, creatorUsername string) error
	// IsUserResponsibleForTender проверяет, имеет ли пользователь права на организацию тендера
	IsUserResponsibleForTender(ctx context.Context, tenderID, username string) error
	// IsUserResponsibleForAuthorBid проверяет, относится ли пользователь к организации автора предложения
	IsUserResponsibleForAuthorBid(ctx context.Context, bidID, username string) error
	// IsUserResponsibleForTenderByBidID проверяет, имеет ли пользователь права на организацию тендера по ID предложения
	IsUserResponsibleForTenderByBidID(ctx context.Context, bidID, username string) error
	// CountResponsibleByBid возвращает количество ответственных пользователей для предложения по его ID
	CountResponsibleByBid(ctx context.Context, bidID string) (int, error)
	// CountApprovedDecisions возвращает количество одобренных решений для предложения по его ID
	CountApprovedDecisions(ctx context.Context, bidID string) (int, error)
}


type Repository interface {
	TenderRepository
	BidRepository
}
