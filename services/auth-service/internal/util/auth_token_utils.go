package internalutil

import (
	"context"
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
)

type AccessClaims struct {
	UserEmail string
	Extra     map[string]any
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	TokenID   string
	UserEmail string
	Extra     map[string]any
	jwt.RegisteredClaims
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func GetRefreshKey(tokenID string) string {
	return fmt.Sprintf("refresh:%s", HashToken(tokenID))
}

func CreateAccessToken(secret []byte, userEmail string, ttl time.Duration, issuer string, extra map[string]any) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserEmail: userEmail,
		Extra:     extra,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func CreateRefreshToken(secret []byte, userEmail string, ttl time.Duration, issuer string, extra map[string]any) (*RefreshClaims, string, error) {
	now := time.Now()
	tokenID := uuid.New().String()

	claims := RefreshClaims{
		TokenID:   tokenID,
		UserEmail: userEmail,
		Extra:     extra,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return nil, "", err
	}

	return &claims, signedToken, nil
}

func ValidateAccessToken(secret []byte, tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func ValidateRefreshToken(secret []byte, tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
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

	accessToken, err := CreateAccessToken([]byte(env.AccessTokenSecret), userEmail, attl, env.AUTH_SRV_NAME, extra)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshClaims, refreshToken, err := CreateRefreshToken([]byte(env.RefreshTokenSecret), userEmail, rttl, env.AUTH_SRV_NAME, extra)
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
