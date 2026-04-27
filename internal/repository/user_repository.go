package repository

import (
	"context"

	"gostudy/internal/model/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository описывает методы для работы с пользователями в хранилище данных.
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}

// userRepository — реализация UserRepository, использующая PostgreSQL.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый экземпляр репозитория пользователей.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// Create сохраняет нового пользователя в базе данных.
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (uuid, name, email, password) VALUES ($1, $2, $3, $4) RETURNING uuid`
	err := r.db.QueryRow(ctx, query, user.Uuid, user.Name, user.Email, user.Password).Scan(&user.Uuid)
	if err != nil {
		return err
	}
	return nil
}

// GetByID возвращает пользователя по его UUID.
func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT uuid, name, email FROM users WHERE uuid = $1`
	var u entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.Uuid, &u.Name, &u.Email)
	return &u, err
}

// GetByEmail возвращает пользователя по его email.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT uuid, name, email, password FROM users WHERE email = $1`
	var u entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.Uuid, &u.Name, &u.Email, &u.Password)
	return &u, err
}

// List возвращает список всех пользователей.
func (r *userRepository) List(ctx context.Context) ([]entity.User, error) {
	query := `SELECT uuid, name, email FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entity.User, 0)

	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.Uuid, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Update обновляет данные пользователя.
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name=$1, email=$2 WHERE uuid=$3`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Uuid)
	return err
}

// Delete удаляет пользователя по его UUID.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE uuid=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
