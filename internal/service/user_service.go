package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/osipovmr/gostudy/internal/db"

	"github.com/osipovmr/gostudy/internal/model/dto"
	"github.com/osipovmr/gostudy/internal/model/entity"
	"github.com/osipovmr/gostudy/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// UserService описывает бизнес-логику для работы с пользователями.
type UserService interface {
	Create(ctx context.Context, input dto.CreateUserInput) (*dto.UserDto, error)
	GetByID(ctx context.Context, id string) (*dto.UserDto, error)
	Update(ctx context.Context, id string, input dto.UpdateUserInput) (*dto.UserDto, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*dto.UserDto, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

// userService — реализация UserService, инкапсулирующая зависимости репозитория и транзакционного менеджера.
type userService struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewUserService создает новый экземпляр сервиса пользователей.
func NewUserService(repo repository.UserRepository, txManager db.TxManager) UserService {
	return &userService{
		userRepository: repo,
		txManager:      txManager,
	}
}

// Create создает нового пользователя, проверяя уникальность email и выполняя операцию в транзакции.
func (s *userService) Create(ctx context.Context, input dto.CreateUserInput) (*dto.UserDto, error) {
	var result *dto.UserDto
	err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		existing, err := s.userRepository.GetByEmail(ctx, input.Email)
		if err == nil && existing != nil {
			return ErrUserExists
		}
		user := &entity.User{
			Uuid:     uuid.New().String(),
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		}
		if err := s.userRepository.Create(ctx, user); err != nil {
			return err
		}
		result = toDTO(user)
		slog.Info("user created", "uuid", result.Uuid)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetByEmail возвращает пользователя по email.
func (s *userService) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetByID возвращает пользователя по UUID.
func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserDto, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return toDTO(user), nil
}

// Update обновляет данные пользователя по UUID.
func (s *userService) Update(ctx context.Context, id string, input dto.UpdateUserInput) (*dto.UserDto, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Name = input.Name
	user.Email = input.Email

	if err := s.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	slog.Info("user updated", "id", user.Uuid)

	return toDTO(user), nil
}

// Delete удаляет пользователя по UUID.
func (s *userService) Delete(ctx context.Context, id string) error {
	_, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	slog.Info("user deleted", "id", id)
	return s.userRepository.Delete(ctx, id)
}

// List возвращает список всех пользователей.
func (s *userService) List(ctx context.Context) ([]*dto.UserDto, error) {
	users, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return toDTOList(users), nil
}

// toDTO преобразует сущность пользователя в DTO.
func toDTO(u *entity.User) *dto.UserDto {
	return &dto.UserDto{
		Uuid:  u.Uuid,
		Name:  u.Name,
		Email: u.Email,
	}
}

// toDTOList преобразует список сущностей пользователей в список DTO.
func toDTOList(users []entity.User) []*dto.UserDto {
	result := make([]*dto.UserDto, 0, len(users))
	for i := range users {
		result = append(result, toDTO(&users[i]))
	}
	return result
}
