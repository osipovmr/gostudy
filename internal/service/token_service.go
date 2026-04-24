package service

import (
	"context"
	"fmt"
	"gostudy/internal/model/entity"
	"gostudy/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type tokenService struct {
	accessSecret    []byte
	refreshSecret   []byte
	accessTTL       time.Duration
	refreshTTL      time.Duration
	tokenRepository repository.TokenRepository
}

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *entity.User) (string, error)
	GenerateRefreshToken(ctx context.Context, user *entity.User) (string, error)
	Save(ctx context.Context, user *entity.User, token string) error
	ValidateAccessToken(ctx context.Context, tokenString string) (*TokenClaims, error)
	ValidateRefreshToken(ctx context.Context, tokenString string) (*TokenClaims, error)
}

func NewTokenService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration, repo repository.TokenRepository) TokenService {
	return &tokenService{
		accessSecret:    []byte(accessSecret),
		refreshSecret:   []byte(refreshSecret),
		accessTTL:       accessTTL,
		refreshTTL:      refreshTTL,
		tokenRepository: repo,
	}
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func (s *tokenService) GenerateAccessToken(ctx context.Context, user *entity.User) (string, error) {
	now := time.Now()

	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Email: user.Email,
		Type:  "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessSecret)
}

func (s *tokenService) GenerateRefreshToken(ctx context.Context, user *entity.User) (string, error) {
	now := time.Now()

	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Email: user.Email,
		Type:  "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecret)
}

func (s *tokenService) Save(ctx context.Context, user *entity.User, token string) error {
	err := s.tokenRepository.Save(ctx, user, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *tokenService) ValidateAccessToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	return s.parseToken(tokenString, s.accessSecret, "access")
}

func (s *tokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	return s.parseToken(tokenString, s.refreshSecret, "refresh")
}

func (s *tokenService) parseToken(tokenString string, secret []byte, expectedType string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.Type != expectedType {
		return nil, fmt.Errorf("invalid token type: %s", claims.Type)
	}

	return claims, nil
}
