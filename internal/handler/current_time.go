package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CurrentTime(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{
		"time": time.Now().UTC().Format(time.RFC3339),
	})
}
