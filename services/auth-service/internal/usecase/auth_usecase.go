package usecase

import (
	"context"
	"fmt"
	"time"

	authservice "github.com/tamirat-dejene/ha-soranu/services/auth-service"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	internalutil "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/util"
	"github.com/tamirat-dejene/ha-soranu/shared/redis"
)

type authUsecase struct {
	ctxTimeout  time.Duration
	userUsecase domain.UserUsecase
	redisClient redis.RedisClient
	userRepo    domain.UserRepository
	environment authservice.Env
}

// LoginWithGoogle implements domain.AuthUsecase.
func (a *authUsecase) LoginWithGoogle(ctx context.Context, google_token string) (domain.UserInfo, domain.AuthToken, error) {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// Verify the Google token and get user info
	userInfo, err := internalutil.VerifyGoogleToken(c, a.environment.GoogleClientID, google_token)
	if err != nil {
		return domain.UserInfo{}, domain.AuthToken{}, err
	}

	email := userInfo.Email

	// Parse token TTLs
	attl, err := time.ParseDuration(a.environment.AccessTokenTTL)
	if err != nil {
		return domain.UserInfo{}, domain.AuthToken{}, err
	}
	rttl, err := time.ParseDuration(a.environment.RefreshTokenTTL)
	if err != nil {
		return domain.UserInfo{}, domain.AuthToken{}, err
	}

	// Generate access token
	accessToken, err := internalutil.CreateToken(
		[]byte(a.environment.AccessTokenSecret),
		email,
		attl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return domain.UserInfo{}, domain.AuthToken{}, err
	}

	// Generate refresh token
	refreshToken, err := internalutil.CreateToken(
		[]byte(a.environment.RefreshTokenSecret),
		email,
		rttl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return domain.UserInfo{}, domain.AuthToken{}, err
	}

	// Store refresh token in Redis
	redisKey := fmt.Sprintf("refresh:%s", refreshToken)
	if err := a.redisClient.Set(redisKey, email, rttl); err != nil {
		return domain.UserInfo{}, domain.AuthToken{}, err
	}

	return userInfo, domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Register implements domain.AuthUsecase.
func (a *authUsecase) Register(ctx context.Context, req domain.CreateUserRequest) (string, domain.AuthToken, error) {
	c, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// Create user using UserRepository
	userID, err := a.userRepo.CreateUser(c, req)
	if err != nil {
		return "", domain.AuthToken{}, err
	}

	// Generate tokens upon successful registration
	attl, err := time.ParseDuration(a.environment.AccessTokenTTL)
	if err != nil {
		return "", domain.AuthToken{}, err
	}
	rttl, err := time.ParseDuration(a.environment.RefreshTokenTTL)
	if err != nil {
		return "", domain.AuthToken{}, err
	}

	accessToken, err := internalutil.CreateToken(
		[]byte(a.environment.AccessTokenSecret),
		req.Email,
		attl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return "", domain.AuthToken{}, err
	}

	refreshToken, err := internalutil.CreateToken(
		[]byte(a.environment.RefreshTokenSecret),
		req.Email,
		rttl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return "", domain.AuthToken{}, err
	}

	// Store refresh token in Redis
	redisKey := fmt.Sprintf("refresh:%s", refreshToken)
	if err := a.redisClient.Set(redisKey, req.Email, rttl); err != nil {
		return "", domain.AuthToken{}, err
	}

	return userID, domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginWithEP implements domain.AuthUsecase.
func (a *authUsecase) LoginWithEP(ctx context.Context, creds domain.LoginCredentials) (domain.AuthToken, error) {
	_, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	email, password := creds.Email, creds.Password

	hashedPassword, err := a.userUsecase.GetUserHashedPassword(ctx, email)
	if err != nil {
		return domain.AuthToken{}, err
	}

	err = internalutil.ComparePassword(hashedPassword, password)
	if err != nil {
		return domain.AuthToken{}, err
	}

	attl, err := time.ParseDuration(a.environment.AccessTokenTTL)
	if err != nil {
		return domain.AuthToken{}, err
	}
	rttl, err := time.ParseDuration(a.environment.RefreshTokenTTL)
	if err != nil {
		return domain.AuthToken{}, err
	}

	// Generate tokens
	accessToken, err := internalutil.CreateToken(
		[]byte(a.environment.AccessTokenSecret),
		email,
		attl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return domain.AuthToken{}, err
	}

	refreshToken, err := internalutil.CreateToken(
		[]byte(a.environment.RefreshTokenSecret),
		email,
		rttl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return domain.AuthToken{}, err
	}

	// Store refresh token in Redis
	redisKey := fmt.Sprintf("refresh:%s", refreshToken)
	err = a.redisClient.Set(redisKey, email, rttl)
	if err != nil {
		return domain.AuthToken{}, err
	}

	return domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout implements domain.AuthUsecase.
func (a *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	_, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	redisKey := fmt.Sprintf("refresh:%s", refreshToken)
	return a.redisClient.Delete(redisKey)
}

// RefreshTokens implements domain.AuthUsecase.
func (a *authUsecase) RefreshTokens(ctx context.Context, refreshToken string) (domain.AuthToken, error) {
	_, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	redisKey := fmt.Sprintf("refresh:%s", refreshToken)

	// Validate refresh token in Redis
	userEmail, err := a.redisClient.Get(redisKey)
	if err != nil {
		return domain.AuthToken{}, fmt.Errorf("invalid or expired refresh token")
	}

	// Optional: delete old refresh token (rotation)
	_ = a.redisClient.Delete(redisKey)

	attl, err := time.ParseDuration(a.environment.AccessTokenTTL)
	if err != nil {
		return domain.AuthToken{}, err
	}

	accessToken, err := internalutil.CreateToken(
		[]byte(a.environment.AccessTokenSecret),
		userEmail,
		attl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return domain.AuthToken{}, err
	}

	rttl, err := time.ParseDuration(a.environment.RefreshTokenTTL)
	if err != nil {
		return domain.AuthToken{}, err
	}

	newRefreshToken, err := internalutil.CreateToken(
		[]byte(a.environment.RefreshTokenSecret),
		userEmail,
		rttl,
		a.environment.AUTH_SRV_NAME,
		nil,
	)
	if err != nil {
		return domain.AuthToken{}, err
	}

	// Store new refresh token in Redis
	newRedisKey := fmt.Sprintf("refresh:%s", newRefreshToken)
	err = a.redisClient.Set(newRedisKey, userEmail, rttl)
	if err != nil {
		return domain.AuthToken{}, err
	}

	return domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// NewAuthUsecase constructor
func NewAuthUsecase(
    userRepo domain.UserRepository,
    userUsecase domain.UserUsecase,
    redisClient redis.RedisClient,
    timeout time.Duration,
    env authservice.Env,
) domain.AuthUsecase {
    return &authUsecase{
        ctxTimeout:  timeout,
        userUsecase: userUsecase,
        redisClient: redisClient,
        userRepo:    userRepo,
        environment: env,
    }
}
