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

type UserInfo struct {
	Username string
	Email    string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

/* ----- Auth Usecase ----- */

type AuthUsecase interface {
	Register(ctx context.Context, req CreateUserRequest) (string, AuthToken, error)
	LoginWithEP(ctx context.Context, creds LoginCredentials) (AuthToken, error)
	LoginWithGoogle(ctx context.Context, googleToken string) (UserInfo, AuthToken, error)
	Logout(ctx context.Context, accessToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (AuthToken, error)
}

/* ----- Auth Repository ----- */

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, userID string, refreshToken string) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error)
}