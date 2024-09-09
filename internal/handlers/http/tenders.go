package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTenders(c *gin.Context) {
	h.log.Info("GetTenders endpoint called")
	c.JSON(http.StatusOK, []string{})
}

func (h *Handler) CreateTender(c *gin.Context) {
	h.log.Info("CreateTender endpoint called")
	c.JSON(http.StatusCreated, gin.H{"message": "Tender created"})
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
