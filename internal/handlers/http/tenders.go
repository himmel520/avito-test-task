package httphandler

import (
	"errors"
	"fmt"
	"net/http"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTenders(c *gin.Context) {
	h.log.Info("GetTenders endpoint called")
	c.JSON(http.StatusOK, []string{})
}

// Создание нового тендера
func (h *Handler) CreateTender(c *gin.Context) {
	var createTender *models.TenderCreate
	if err := c.BindJSON(&createTender); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{fmt.Sprintf("ошибка в теле запроса: %v", err)})
		return
	}

	tender, err := h.srv.CreateTender(c.Request.Context(), createTender)
	switch {
	case errors.Is(err, repository.ErrUserNotExist):
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrRelationNotExist):
		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{err.Error()})
		return
	case errors.Is(err, repository.ErrOrganizationDepencyNotFound):
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{err.Error()})
		return
	case err != nil:
		h.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, tender)
}

func (h *Handler) GetMyTenders(c *gin.Context) {
	h.log.Info("GetMyTenders endpoint called")
	c.JSON(http.StatusOK, []string{})
}

func (h *Handler) GetTenderStatus(c *gin.Context) {
	h.log.Info("GetTenderStatus endpoint called")
	c.JSON(http.StatusOK, gin.H{"status": "status"})
}

func (h *Handler) UpdateTenderStatus(c *gin.Context) {
	h.log.Info("UpdateTenderStatus endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Tender status updated"})
}

func (h *Handler) EditTender(c *gin.Context) {
	h.log.Info("EditTender endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Tender edited"})
}

func (h *Handler) RollbackTenderVersion(c *gin.Context) {
	h.log.Info("RollbackTenderVersion endpoint called")
	c.JSON(http.StatusOK, gin.H{"message": "Tender version rolled back"})
}
