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

// AddAddress implements domain.UserUseCase.
func (u *userUsecase) AddAddress(ctx context.Context, userID string, address domain.Address) (*domain.Address, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	
	addr, err := u.userRepository.AddAddress(c, userID, address)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

// AddPhoneNumber implements domain.UserUseCase.
func (u *userUsecase) AddPhoneNumber(ctx context.Context, userID string, phone string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.AddPhoneNumber(c, userID, phone)
}

// GetAddresses implements domain.UserUseCase.
func (u *userUsecase) GetAddresses(ctx context.Context, userID string) ([]domain.Address, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.GetAddresses(c, userID)
}

// GetUser implements domain.UserUseCase.
func (u *userUsecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.GetUserByID(c, userID)
}

// RemoveAddress implements domain.UserUseCase.
func (u *userUsecase) RemoveAddress(ctx context.Context, userID string, addressID string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.RemoveAddress(c, userID, addressID)
}

// RemovePhoneNumber implements domain.UserUseCase.
func (u *userUsecase) RemovePhoneNumber(ctx context.Context, userID string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.RemovePhoneNumber(c, userID)
}

// UpdatePhoneNumber implements domain.UserUseCase.
func (u *userUsecase) UpdatePhoneNumber(ctx context.Context, userID string, phone string) error {
	c, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.userRepository.UpdatePhoneNumber(c, userID, phone)
}

func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &userUsecase{
		userRepository: userRepo,
		ctxTimeout:     timeout,
	}
}
