package usecase

import (
	"context"
	"time"

	authservice "github.com/tamirat-dejene/ha-soranu/services/auth-service"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	internalutil "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/util"
	jwtvalidator "github.com/tamirat-dejene/ha-soranu/shared/pkg/auth/jwtvalidator"
)

type authUsecase struct {
	ctxTimeout time.Duration
	authRepo   domain.AuthRepository
	userRepo   domain.UserRepository
	env        authservice.Env
}

// LoginWithEmailAndPassword implements domain.AuthUseCase.
func (a *authUsecase) LoginWithEmailAndPassword(ctx context.Context, input domain.LoginWithEmail) (*domain.User, *domain.AuthTokens, error) {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	passwordHash, err := a.userRepo.GetUserPasswordHashByEmail(c, input.Email)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	err = internalutil.ComparePassword(passwordHash, input.Password)
	if err != nil {
		return nil, nil, errs.ErrInvalidCredentials
	}

	authToken, refreshTokenID, err := internalutil.SignUser(input.Email, &a.env, nil)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	err = a.authRepo.SaveRefreshToken(c, input.Email, refreshTokenID)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	user, err := a.userRepo.GetUserByEmail(c, input.Email)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	return user, authToken, nil
}

func (a *authUsecase) LoginWithGoogle(ctx context.Context, input domain.LoginWithGoogle) (*domain.User, *domain.AuthTokens, error) {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// 1. Validate Google token
	claims, err := internalutil.ValidateGoogleIDToken(c, a.env.GoogleClientID, input.IDToken)
	if err != nil {
		return nil, nil, err
	}

	// 2. Try to find the user
	user, err := a.userRepo.GetUserByEmail(c, claims.Email)

	// 3. If user does not exist → create a new user
	if err == nil && user == nil {
		newUser := &domain.UserRegister{
			Email:       claims.Email,
			Username:    claims.Username,
			PhoneNumber: "",
			Password:    "",
		}

		user, err = a.userRepo.CreateUser(c, newUser)
		if err != nil {
			return nil, nil, err
		}
	} else {
		// other DB errors
		return nil, nil, err
	}

	// 4. Generate JWT and refresh token
	authToken, refreshTokenID, err := internalutil.SignUser(user.Email, &a.env, nil)
	if err != nil {
		return nil, nil, err
	}

	// 5. Store refresh token
	err = a.authRepo.SaveRefreshToken(c, user.UserID, refreshTokenID)
	if err != nil {
		return nil, nil, err
	}

	// 6. Return user + tokens
	return user, authToken, nil
}

// Logout implements domain.AuthUseCase.
func (a *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	claims, err := jwtvalidator.ValidateRefreshToken(a.env.RefreshTokenPublicKey, refreshToken)
	if err != nil {
		return errs.ErrTokenRevoked
	}

	return a.authRepo.DeleteRefreshToken(c, claims.TokenID)
}

// RefreshTokens implements domain.AuthUseCase.
func (a *authUsecase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthTokens, error) {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	claims, err := jwtvalidator.ValidateRefreshToken(a.env.RefreshTokenPublicKey, refreshToken)
	if err != nil {
		return nil, errs.ErrTokenRevoked
	}
	_, err = a.authRepo.ConsumeRefreshToken(c, claims.TokenID)
	if err != nil {
		return nil, errs.ErrTokenRevoked
	}

	authToken, _, err := internalutil.SignUser(claims.UserEmail, &a.env, nil)
	if err != nil {
		return nil, errs.ErrInternalServer
	}

	return authToken, nil
}

// Register implements domain.AuthUseCase.
func (a *authUsecase) Register(ctx context.Context, input domain.UserRegister) (*domain.User, *domain.AuthTokens, error) {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var err error
	input.Password, err = internalutil.HashPassword(input.Password)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	user, err := a.userRepo.CreateUser(c, &input)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}
	authToken, refreshTokenID, err := internalutil.SignUser(user.Email, &a.env, nil)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	err = a.authRepo.SaveRefreshToken(c, user.UserID, refreshTokenID)
	if err != nil {
		return nil, nil, errs.ErrInternalServer
	}

	return user, authToken, nil
}

// NewAuthUsecase constructor
func NewAuthUsecase(ctxTimeout time.Duration, authRepo domain.AuthRepository, userRepo domain.UserRepository, env authservice.Env) domain.AuthUseCase {
	return &authUsecase{
		ctxTimeout: ctxTimeout,
		authRepo:   authRepo,
		userRepo:   userRepo,
		env:        env,
	}
}
