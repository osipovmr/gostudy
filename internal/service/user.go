package service

import (
	"gostudy/internal/model"
	"gostudy/internal/repository"
)

type UserService interface {
	Create(user *model.User) error
	GetByID(id string) (*model.User, error)
	Update(id string, user *model.User) error
	Delete(id string) error
	List() ([]*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(user *model.User) error {
	return s.repo.Create(user)
}

func (s *userService) GetByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) Update(id string, user *model.User) error {
	return s.repo.Update(id, user)
}

func (s *userService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *userService) List() ([]*model.User, error) {
	return s.repo.List()
}
