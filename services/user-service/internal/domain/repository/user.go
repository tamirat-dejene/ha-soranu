package repository

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/user-service/internal/domain/entity"
	"github.com/tamirat-dejene/ha-soranu/user-service/internal/domain/errs"
)

type inMemoryUserRepository struct {
	users map[string]*entity.User
}

func NewInMemoryUserRepository() entity.UserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *inMemoryUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if _, exists := r.users[user.ID]; exists {
		return errs.ErrUserAlreadyExists
	}
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *inMemoryUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return errs.ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepository) DeleteUser(ctx context.Context, id string) error {
	if _, exists := r.users[id]; !exists {
		return errs.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}
