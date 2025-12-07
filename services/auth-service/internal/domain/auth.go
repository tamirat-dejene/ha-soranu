package domain

import "context"


/* ----- Auth Entities ----- */

type LoginCredentials struct {
	Email    string
	Password string
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

/* ----- Auth Usecase ----- */

type AuthUsecase interface {
	Register(ctx context.Context, req CreateUserRequest) (string, AuthToken, error)
	Login(ctx context.Context, creds LoginCredentials) (AuthToken, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (string, error)
}