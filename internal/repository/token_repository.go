package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}
