package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	healthResponse := map[string]string{"status": "ok"}
	c.JSON(http.StatusOK, healthResponse)
}
