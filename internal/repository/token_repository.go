package repository

import (
	"context"
	"gostudy/internal/model/entity"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	Save(ctx context.Context, user *entity.User, token string) error
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Save(ctx context.Context, user *entity.User, token string) error {
	query := `INSERT INTO refresh_token (user_uuid, token_hash, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.Uuid, token, time.Now())
	if err != nil {
		return err
	}
	return nil
}
