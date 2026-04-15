package repository

import (
	"gostudy/internal/model"
	"sync"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id string) (*model.User, error)
	Update(id string, user *model.User) error
	Delete(id string) error
	List() ([]*model.User, error)
}

type userRepo struct {
	store map[string]*model.User
	mu    sync.RWMutex
}

func NewUserRepository() UserRepository {
	return &userRepo{
		store: make(map[string]*model.User),
		mu:    sync.RWMutex{},
	}
}

func (r *userRepo) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[user.ID]; exists {
		return model.ErrUserExists
	}
	r.store[user.ID] = user
	return nil
}

func (r *userRepo) GetByID(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.store[id]
	if !exists {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}

func (r *userRepo) Update(id string, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return model.ErrUserNotFound
	}
	user.ID = id
	r.store[id] = user
	return nil
}

func (r *userRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return model.ErrUserNotFound
	}
	delete(r.store, id)
	return nil
}

func (r *userRepo) List() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.store))
	for _, user := range r.store {
		users = append(users, user)
	}
	return users, nil
}
