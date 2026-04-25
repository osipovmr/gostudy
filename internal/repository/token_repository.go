package repository

import (
	"context"
	"fmt"
	"gostudy/internal/model/entity"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	Save(ctx context.Context, user *entity.User, token string) error
	Delete(ctx context.Context, user *entity.User) error
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

func (r *tokenRepository) Delete(ctx context.Context, user *entity.User) error {
	query := `DELETE FROM refresh_token WHERE user_uuid = $1`

	tag, err := r.db.Exec(ctx, query, user.Uuid)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("refresh token not found")
	}
	return nil
}
