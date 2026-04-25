package middleware

import (
	"gostudy/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	tokenService service.TokenService
}

func NewAuthMiddleware(s service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: s}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := m.tokenService.ValidateAccessToken(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userEmail", claims.Email)
		c.Set("accessToken", parts[1])
		c.Next()
	}
}
