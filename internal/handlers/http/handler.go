package httphandler

import (
	service "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service interface {
	service.TenderService
	service.BidService
}

type Handler struct {
	srv Service
	log *logrus.Logger
}

func New(srv Service, log *logrus.Logger) *Handler {
	return &Handler{srv: srv, log: log}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/ping", h.Ping) // Проверка доступности сервера

		tenders := api.Group("/tenders")
		{
			tenders.GET("", h.GetTenders) // Получение списка тендеров
			tenders.POST("/new", h.CreateTender) // Создание нового тендера
			tenders.GET("/my", h.GetMyTenders) // Получить тендеры пользователя
			tenders.GET("/:tenderId/status", h.GetTenderStatus) // Получение текущего статуса тендера
			tenders.PUT("/:tenderId/status", h.UpdateTenderStatus) // Изменение статуса тендера
			tenders.PATCH("/:tenderId/edit", h.EditTender) // Редактирование тендера
			tenders.PUT("/:tenderId/rollback/:version", h.RollbackTenderVersion) // Откат версии тендера
		}

		bids := api.Group("/bids")
		{
			bids.POST("/new", h.CreateBid) // Создание нового предложения
			bids.GET("/my", h.GetMyBids) // Получение списка ваших предложений
			bids.GET("/:tenderId/list", h.GetBidsForTender) // Получение списка предложений для тендера
			bids.GET("/:bidId/status", h.GetBidStatus) // Получение текущего статуса предложения
			bids.PUT("/:bidId/status", h.UpdateBidStatus) // Изменение статуса предложения
			bids.PATCH("/:bidId/edit", h.EditBid) // Редактирование параметров предложения
			bids.PUT("/:bidId/submit_decision", h.SubmitDecision) // Отправка решения по предложению
			bids.PUT("/:bidId/feedback", h.SubmitFeedback) // Отправка отзыва по предложению
			bids.PUT("/:bidId/rollback/:version", h.RollbackBidVersion) // Откат версии предложения
			bids.GET("/:tenderId/reviews", h.GetBidReviews) // Просмотр отзывов на прошлые предложения
		}
	}

	return r
}

