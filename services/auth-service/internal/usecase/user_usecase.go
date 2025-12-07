package usecase

import (
	"context"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
)

type userUsecase struct {
	userRepository domain.UserRepository
	ctxTimeout     time.Duration
}

// AddUserAddress implements domain.UserUsecase.
func (u *userUsecase) AddUserAddress(ctx context.Context, id string, address string) error {
	panic("unimplemented")
}

// DeleteUser implements domain.UserUsecase.
func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetUserAddresses implements domain.UserUsecase.
func (u *userUsecase) GetUserAddresses(ctx context.Context, id string) ([]string, error) {
	panic("unimplemented")
}

// GetUserByID implements domain.UserUsecase.
func (u *userUsecase) GetUserByID(ctx context.Context, id string) (domain.UserResponse, error) {
	panic("unimplemented")
}

// RemoveUserAddress implements domain.UserUsecase.
func (u *userUsecase) RemoveUserAddress(ctx context.Context, id string, address string) error {
	panic("unimplemented")
}

// UpdateUser implements domain.UserUsecase.
func (u *userUsecase) UpdateUser(ctx context.Context, id string, req domain.UpdateUserRequest) (domain.UserResponse, error) {
	panic("unimplemented")
}

func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepo,
		ctxTimeout:     timeout,
	}
}
