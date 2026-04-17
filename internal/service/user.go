package service

import (
	"context"
	"errors"
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
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, input dto.CreateUserInput) (*dto.UserDto, error) {
	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, ErrUserExists
	}
	user := &entity.User{
		ID:    uuid.New().String(),
		Name:  input.Name,
		Email: input.Email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	slog.Info("user created", "id", user.ID)
	return toDTO(user), nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.UserDto, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return toDTO(user), nil
}

func (s *userService) Update(ctx context.Context, id string, input dto.UpdateUserInput) (*dto.UserDto, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Name = input.Name
	user.Email = input.Email

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	slog.Info("user updated", "id", user.ID)

	return toDTO(user), nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}
	slog.Info("user deleted", "id", id)
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*dto.UserDto, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return toDTOList(users), nil
}

func toDTO(u *entity.User) *dto.UserDto {
	return &dto.UserDto{
		ID:    u.ID,
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
