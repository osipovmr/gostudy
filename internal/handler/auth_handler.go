package handler

import (
	"gostudy/internal/facade"
	"gostudy/internal/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authFacade facade.AuthFacade
}

func NewAuthHandler(authFacade facade.AuthFacade) *AuthHandler {
	return &AuthHandler{authFacade: authFacade}
}

func (h *AuthHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	//routerGroup.GET("/me", h.GetMe)
	routerGroup.POST("/register", h.Register)
	routerGroup.POST("/login", h.Login)
	//routerGroup.POST("/refresh", h.Refresh)
	//.POST("/logout", h.Logout)
}

//func (h *AuthHandler) GetMe(c *gin.Context) {
//	// Пример: если у тебя userID кладется в middleware в context
//	userID, ok := c.Get("userID")
//	if !ok {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//		return
//	}
//
//	user, err := h.authFacade.GetMe(c.Request.Context(), userID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, user)
//}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authFacade.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authFacade.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

//func (h *AuthHandler) Refresh(c *gin.Context) {
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
//		return
//	}
//
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
//		return
//	}
//
//	token := parts[1]
//
//	res, err := h.authFacade.Refresh(c.Request.Context(), token)
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, res)
//}
//
//func (h *AuthHandler) Logout(c *gin.Context) {
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
//		return
//	}
//
//	parts := strings.SplitN(authHeader, " ", 2)
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
//		return
//	}
//
//	token := parts[1]
//
//	if err := h.authFacade.Logout(c.Request.Context(), token); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.Status(http.StatusNoContent)
//}
