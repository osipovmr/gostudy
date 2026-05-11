package handler

import (
	"log/slog"
	"net/http"

	"github.com/osipovmr/gostudy/internal/facade"
	"github.com/osipovmr/gostudy/internal/model/dto"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authFacade facade.AuthFacade
}

func NewAuthHandler(authFacade facade.AuthFacade) *AuthHandler {
	return &AuthHandler{authFacade: authFacade}
}

func (h *AuthHandler) RegisterRoutes(routerGroup *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	public := routerGroup.Group("")
	{
		public.POST("/register", h.Register)
		public.POST("/login", h.Login)
		public.POST("/refresh", h.Refresh)
	}

	protected := routerGroup.Group("")
	protected.Use(authMiddleware)
	{
		protected.GET("/me", h.GetMe)
		protected.POST("/logout", h.Logout)
	}
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userEmail, _ := c.Get("userEmail")
	user, err := h.authFacade.GetMe(c.Request.Context(), userEmail.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.authFacade.Refresh(c.Request.Context(), req)
	if err != nil {
		slog.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userEmail, _ := c.Get("userEmail")
	if err := h.authFacade.Logout(c.Request.Context(), userEmail.(string)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}
