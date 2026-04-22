package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) TxManager {
	return &txManager{pool: pool}
}

func (m *txManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, txKey{}, tx)
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err := fn(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
