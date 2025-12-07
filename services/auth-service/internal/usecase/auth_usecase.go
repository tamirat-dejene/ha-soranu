package usecase

import (
	"context"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
)

type authUsecase struct {
	ctxTimeout time.Duration
}

// Login implements domain.AuthUsecase.
func (a *authUsecase) Login(ctx context.Context, creds domain.LoginCredentials) (domain.AuthToken, error) {
	panic("unimplemented")
}

// Logout implements domain.AuthUsecase.
func (a *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	panic("unimplemented")
}

// RefreshTokens implements domain.AuthUsecase.
func (a *authUsecase) RefreshTokens(ctx context.Context, refreshToken string) (string, error) {
	panic("unimplemented")
}

// Register implements domain.AuthUsecase.
func (a *authUsecase) Register(ctx context.Context, req domain.CreateUserRequest) (string, domain.AuthToken, error) {
	panic("unimplemented")
}

func NewAuthUsecase(timeout time.Duration) domain.AuthUsecase {
	return &authUsecase{
		ctxTimeout: timeout,
	}
}
