package httphandler

import (
	"errors"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/gin-gonic/gin"
)

// Создание нового предложения
func (h *Handler) CreateBid(c *gin.Context) {
	var createBid *models.BidCreate
	if err := c.BindJSON(&createBid); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("ошибка в теле запроса: %v", err)})
		return
	}

	bid, err := h.srv.CreateBid(c.Request.Context(), createBid)
	switch {
	case errors.Is(err, repository.ErrBidUnique) || errors.Is(err, repository.ErrTenderClosed):
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidDependencyNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// Получение списка предложений пользователя
func (h *Handler) GetMyBids(c *gin.Context) {
	var query myQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	bids, err := h.srv.GetUserBids(c.Request.Context(), query.Username, query.Limit, query.Offset)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// Получение списка предложений для тендера
func (h *Handler) GetBidsForTender(c *gin.Context) {
	var uri bidTenderIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query myQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	bids, err := h.srv.GetBidsForTender(c.Request.Context(), uri.ID, query.Username, query.Limit, query.Offset)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidTenderNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// Получение текущего статуса предложения
func (h *Handler) GetBidStatus(c *gin.Context) {
	var uri bidIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query UsernameQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	status, err := h.srv.GetBidStatus(c.Request.Context(), uri.ID, query.Username)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// Изменение статуса предложения
func (h *Handler) UpdateBidStatus(c *gin.Context) {
	var uri bidIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query updateBidStatusQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	bid, err := h.srv.UpdateBidStatus(c.Request.Context(), uri.ID, query.Username, &query.Status)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// Редактирование параметров предложения
func (h *Handler) EditBid(c *gin.Context) {
	var uri bidIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query UsernameQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	var bidEdit *models.BidEdit
	if err := c.BindJSON(&bidEdit); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("ошибка в теле запроса: %v", err)})
		return
	}

	if bidEdit.IsEmpty() {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"нет изменений"})
		return
	}

	bid, err := h.srv.EditBid(c.Request.Context(), uri.ID, query.Username, bidEdit)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// Отправка решения по предложению
func (h *Handler) SubmitDecision(c *gin.Context) {
	var uri bidIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query decisionQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	bid, err := h.srv.SubmitBidDecision(c.Request.Context(), uri.ID, query.Username, &query.Decision)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidNotFound) || errors.Is(err, repository.ErrTenderNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// Отправка отзыва по предложению
func (h *Handler) SubmitFeedback(c *gin.Context) {
	var uri bidIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
	}

	var query feedbackQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	bid, err := h.srv.SubmitBidFeedback(c.Request.Context(), uri.ID, query.Username, &query.BidFeedback)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// Откат версии предложения
func (h *Handler) RollbackBidVersion(c *gin.Context) {
	var uri rollbackBidUri
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query UsernameQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	bid, err := h.srv.RollbackBid(c.Request.Context(), uri.ID, query.Username, uri.Version)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidORVersionNotFound) || errors.Is(err, repository.ErrBidNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// Просмотр отзывов на прошлые предложения
func (h *Handler) GetBidReviews(c *gin.Context) {
	var uri bidIdURI
	if err := c.BindUri(&uri); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный uri: %v", err)})
		return
	}

	var query reviewsQuery
	if err := c.BindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("неккоректный query: %v", err)})
		return
	}

	feedbacks, err := h.srv.GetBidReviews(c.Request.Context(), uri.ID, query.AuthorUsername, query.RequesterUsername, query.Limit, query.Offset)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrBidReviewsNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	c.JSON(http.StatusOK, feedbacks)
}
