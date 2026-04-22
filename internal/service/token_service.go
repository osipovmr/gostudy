package service

import (
	"gostudy/internal/repository"
)

type TokenService interface {
}
type tokenService struct {
	tokenRepository repository.TokenRepository
}

func NewTokenService(repo repository.TokenRepository) TokenService {
	return &tokenService{
		tokenRepository: repo,
	}
}
