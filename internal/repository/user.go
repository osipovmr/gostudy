package repository

import (
	"context"
	"gostudy/internal/model/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	List(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, name) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, user.ID, user.Name).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, name FROM users WHERE id=$1`

	var u entity.User
	if err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) List(ctx context.Context) ([]entity.User, error) {
	query := `SELECT id, name FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, user.Name, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
