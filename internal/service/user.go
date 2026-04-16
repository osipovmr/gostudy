package service

import (
	"context"
	"gostudy/internal/model/entity"
	"gostudy/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]entity.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, user *entity.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, user *entity.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]entity.User, error) {
	return s.repo.List(ctx)
}
