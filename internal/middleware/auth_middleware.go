package middleware

import (
	"net/http"
	"strings"

	"github.com/osipovmr/gostudy/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// UserClaims представляет пользовательские JWT-claims, используемые в приложении.
type UserClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware отвечает за проверку JWT access-токена в HTTP-запросах.
type AuthMiddleware struct {
	tokenService service.TokenService
}

// NewAuthMiddleware создает новый экземпляр middleware для аутентификации.
func NewAuthMiddleware(s service.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: s}
}

// Handler проверяет наличие и валидность Bearer access-токена и добавляет данные пользователя в контекст запроса.
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
		// сохраняем access token в контексте запроса
		c.Set("accessToken", parts[1])
		c.Next()
	}
}
