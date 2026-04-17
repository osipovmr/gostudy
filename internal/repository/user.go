package repository

import (
	"context"

	"gostudy/internal/model/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
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
	query := `INSERT INTO "user" (id, name, email) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT id, name, email FROM "user" WHERE id = $1`
	var u entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email)
	return &u, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, name, email FROM "user" WHERE email = $1`
	var u entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email)
	return &u, err
}

func (r *userRepository) List(ctx context.Context) ([]entity.User, error) {
	query := `SELECT id, name, email FROM "user"`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entity.User, 0)

	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE "user" SET name=$1, email=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "user" WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
