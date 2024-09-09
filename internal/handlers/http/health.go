package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Проверка доступности сервера
func(h *Handler) Ping(c *gin.Context){
	c.String(http.StatusOK, "ok")
}
