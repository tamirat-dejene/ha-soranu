package internalutil

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authservice "github.com/tamirat-dejene/ha-soranu/services/auth-service"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	jwtvalidator "github.com/tamirat-dejene/ha-soranu/shared/pkg/auth/jwtvalidator"
)

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func GetRefreshKey(tokenID string) string {
	return fmt.Sprintf("refresh:%s", HashToken(tokenID))
}

func CreateAccessToken(privateKey *rsa.PrivateKey, userEmail string, ttl time.Duration, issuer string, extra map[string]any) (string, error) {
	now := time.Now()
	claims := jwtvalidator.AccessClaims{
		UserEmail: userEmail,
		Extra:     extra,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func CreateRefreshToken(privateKey *rsa.PrivateKey, userEmail string, ttl time.Duration, issuer string, extra map[string]any) (*jwtvalidator.RefreshClaims, string, error) {
	now := time.Now()
	tokenID := uuid.New().String()

	claims := jwtvalidator.RefreshClaims{
		TokenID:   tokenID,
		UserEmail: userEmail,
		Extra:     extra,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, "", err
	}

	return &claims, signedToken, nil
}

func SignUser(userEmail string, env *authservice.Env, extra map[string]any) (*domain.AuthTokens, string, error) {
	attl, err := time.ParseDuration(env.AccessTokenTTL)
	if err != nil {
		return nil, "", fmt.Errorf("invalid access token TTL: %w", err)
	}

	rttl, err := time.ParseDuration(env.RefreshTokenTTL)
	if err != nil {
		return nil, "", fmt.Errorf("invalid refresh token TTL: %w", err)
	}

	accessPrivateKey, err := jwtvalidator.ParseRSAPrivateKeyFromString(env.AccessTokenPrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("invalid access token private key: %w", err)
	}

	refreshPrivateKey, err := jwtvalidator.ParseRSAPrivateKeyFromString(env.RefreshTokenPrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("invalid refresh token private key: %w", err)
	}

	accessToken, err := CreateAccessToken(accessPrivateKey, userEmail, attl, env.AUTH_SRV_NAME, extra)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshClaims, refreshToken, err := CreateRefreshToken(refreshPrivateKey, userEmail, rttl, env.AUTH_SRV_NAME, extra)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, refreshClaims.TokenID, nil
}

func ValidateGoogleIDToken(ctx context.Context, client_id, id_token string) (domain.User, error) {
	payload, err := idtoken.Validate(ctx, id_token, client_id)
	if err != nil {
		return domain.User{}, errors.New("invalid Google token")
	}

	// Extract email and name from payload
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return domain.User{}, errors.New("email not found in token")
	}

	name, _ := payload.Claims["name"].(string)

	user := domain.User{
		Email:    email,
		Username: name,
	}

	return user, nil
}
