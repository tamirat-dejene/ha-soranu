package domain

import "context"

type AuthTokens struct {
    AccessToken  string
    RefreshToken string
}

type UserRegister struct {
    Email       string
    Password    string
    Username    string
    PhoneNumber string
}

type LoginWithEmail struct {
    Email    string
    Password string
}

type LoginWithGoogle struct {
    IDToken string
}

type AuthUseCase interface {
    Register(ctx context.Context, input UserRegister) (*User, *AuthTokens, error)
    LoginWithEmailAndPassword(ctx context.Context, input LoginWithEmail) (*User, *AuthTokens, error)
    LoginWithGoogle(ctx context.Context, input LoginWithGoogle) (*User, *AuthTokens, error)
    Logout(ctx context.Context, refreshToken string) error
    RefreshTokens(ctx context.Context, refreshToken string) (*AuthTokens, error)
}

type AuthRepository interface {
    SaveRefreshToken(ctx context.Context, userID string, tokenId string) error
    DeleteRefreshToken(ctx context.Context, tokenId string) error
    ValidateRefreshToken(ctx context.Context, tokenId string) (string, error)
    ConsumeRefreshToken(ctx context.Context, tokenId string) (string, error)
}