package httphandler

import (
	"errors"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBid(c *gin.Context) {
	var createBid *models.BidCreate
	if err := c.BindJSON(&createBid); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("ошибка в теле запроса: %v", err)})
		return
	}

	bid, err := h.srv.CreateBid(c.Request.Context(), createBid)
	switch {
	case errors.Is(err, repository.ErrBidUnique):
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

func (h *Handler) GetMyBids(c *gin.Context) {
	h.log.Info("GetMyBids endpoint called")
	c.JSON(http.StatusOK, []string{})
}

func (h *Handler) GetBidsForTender(c *gin.Context) {
	h.log.Info("GetBidsForTender endpoint called")
	c.JSON(http.StatusOK, []string{})
}

func (h *Handler) GetBidStatus(c *gin.Context) {
	h.log.Info("GetBidStatus endpoint called")
	c.JSON(http.StatusOK, gin.H{"status": "status"})
}

func (h *Handler) UpdateBidStatus(c *gin.Context) {
	h.log.Info("UpdateBidStatus endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Bid status updated"})
}

func (h *Handler) EditBid(c *gin.Context) {
	h.log.Info("EditBid endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Bid edited"})
}

func (h *Handler) SubmitDecision(c *gin.Context) {
	h.log.Info("SubmitDecision endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Decision submitted"})
}

func (h *Handler) SubmitFeedback(c *gin.Context) {
	h.log.Info("SubmitFeedback endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Feedback submitted"})
}

func (h *Handler) RollbackBidVersion(c *gin.Context) {
	h.log.Info("RollbackBidVersion endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Bid version rolled back"})
}

func (h *Handler) GetBidReviews(c *gin.Context) {
	h.log.Info("GetBidReviews endpoint called")
	c.JSON(http.StatusOK, []string{})
}
