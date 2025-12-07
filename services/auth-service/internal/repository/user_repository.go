package repository

import (
	"context"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
)

type userRepository struct {
	// db connection or any other dependencies
}

// AddAddress implements domain.UserRepository.
func (u *userRepository) AddAddress(ctx context.Context, id string, address string) error {
	panic("unimplemented")
}

// Delete implements domain.UserRepository.
func (u *userRepository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindByID implements domain.UserRepository.
func (u *userRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	panic("unimplemented")
}

// GetAddresses implements domain.UserRepository.
func (u *userRepository) GetAddresses(ctx context.Context, id string) ([]string, error) {
	panic("unimplemented")
}

// RemoveAddress implements domain.UserRepository.
func (u *userRepository) RemoveAddress(ctx context.Context, id string, address string) error {
	panic("unimplemented")
}

// Update implements domain.UserRepository.
func (u *userRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	panic("unimplemented")
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() domain.UserRepository {
	return &userRepository{}
}
