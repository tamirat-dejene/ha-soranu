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

// CreateUser implements domain.UserUsecase.
func (u *userUsecase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (string, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.CreateUser(c, req)
}

// GetUserHashedPassword implements domain.UserUsecase.
func (u *userUsecase) GetUserHashedPassword(ctx context.Context, email string) (string, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.GetHashedPassword(c, email)
}

// AddUserAddress implements domain.UserUsecase.
func (u *userUsecase) AddUserAddress(ctx context.Context, id string, address string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.AddAddress(c, id, address)
}

// DeleteUser implements domain.UserUsecase.
func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.Delete(c, id)
}

// GetUserAddresses implements domain.UserUsecase.
func (u *userUsecase) GetUserAddresses(ctx context.Context, id string) ([]string, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.GetAddresses(c, id)
}

// GetUserByID implements domain.UserUsecase.
func (u *userUsecase) GetUserByID(ctx context.Context, id string) (domain.UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	user, err := u.userRepository.FindByID(c, id)
	if err != nil {
		return domain.UserResponse{}, err
	}

	return domain.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Addressess: user.Addressess,
	}, nil
}

// RemoveUserAddress implements domain.UserUsecase.
func (u *userUsecase) RemoveUserAddress(ctx context.Context, id string, address string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.RemoveAddress(c, id, address)
}

// UpdateUser implements domain.UserUsecase.
func (u *userUsecase) UpdateUser(ctx context.Context, id string, req domain.UpdateUserRequest) (domain.UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	user, err := u.userRepository.FindByID(c, id)
	if err != nil {
		return domain.UserResponse{}, err
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		user.Password = *req.Password
	}
	if req.Addressess != nil {
		user.Addressess = *req.Addressess
	}

	updatedUser, err := u.userRepository.Update(c, user)
	if err != nil {
		return domain.UserResponse{}, err
	}

	return domain.UserResponse{
		ID:         updatedUser.ID,
		Username:   updatedUser.Username,
		Email:      updatedUser.Email,
		Addressess: updatedUser.Addressess,
	}, nil
}

func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepo,
		ctxTimeout:     timeout,
	}
}
