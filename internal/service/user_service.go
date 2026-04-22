package service

import (
	"context"
	"errors"
	"gostudy/internal/db"
	"log/slog"

	"gostudy/internal/model/dto"
	"gostudy/internal/model/entity"
	"gostudy/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type UserService interface {
	Create(ctx context.Context, input dto.CreateUserInput) (*dto.UserDto, error)
	GetByID(ctx context.Context, id string) (*dto.UserDto, error)
	Update(ctx context.Context, id string, input dto.UpdateUserInput) (*dto.UserDto, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*dto.UserDto, error)
}

type userService struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewUserService(repo repository.UserRepository, txManager db.TxManager) UserService {
	return &userService{
		userRepository: repo,
		txManager:      txManager,
	}
}

func (s *userService) Create(ctx context.Context, input dto.CreateUserInput) (*dto.UserDto, error) {
	var result *dto.UserDto
	err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		existing, err := s.userRepository.GetByEmail(ctx, input.Email)
		if err == nil && existing != nil {
			return ErrUserExists
		}
		user := &entity.User{
			Uuid:  uuid.New().String(),
			Name:  input.Name,
			Email: input.Email,
			//todo хешировать
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

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserDto, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return toDTO(user), nil
}

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

func (s *userService) Delete(ctx context.Context, id string) error {
	_, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	slog.Info("user deleted", "id", id)
	return s.userRepository.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*dto.UserDto, error) {
	users, err := s.userRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	return toDTOList(users), nil
}

func toDTO(u *entity.User) *dto.UserDto {
	return &dto.UserDto{
		Uuid:  u.Uuid,
		Name:  u.Name,
		Email: u.Email,
	}
}

func toDTOList(users []entity.User) []*dto.UserDto {
	result := make([]*dto.UserDto, 0, len(users))
	for i := range users {
		result = append(result, toDTO(&users[i]))
	}
	return result
}
