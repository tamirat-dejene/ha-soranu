package internalutil

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
)

type CustomClaims struct {
	UserEmail string         `json:"userEmail"`
	Extra     map[string]any `json:"extra,omitempty"`
	jwt.RegisteredClaims
}

type TokenClaims struct {
	Sub   string         `json:"sub"`
	Iat   int64          `json:"iat"`
	Exp   int64          `json:"exp"`
	Iss   string         `json:"iss,omitempty"`
	Extra map[string]any `json:"extra,omitempty"`
}

// CreateToken generates a JWT token with the given parameters.
func CreateToken(secret []byte, userEmail string, ttl time.Duration, issuer string, extra map[string]any) (string, error) {
	now := time.Now()

	claims := CustomClaims{
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

// ValidateToken validates the given JWT token string and returns the claims if valid.
func ValidateToken(secret []byte, tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		// ensure alg=HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// VerifyGoogleToken verifies a Google ID token and returns user info.
func VerifyGoogleToken(ctx context.Context, client_id, google_token string) (domain.UserInfo, error) {
	payload, err := idtoken.Validate(ctx, google_token, client_id)
	if err != nil {
		return domain.UserInfo{}, errors.New("invalid Google token")
	}

	// Extract email and name from payload
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return domain.UserInfo{}, errors.New("email not found in token")
	}

	name, _ := payload.Claims["name"].(string)

	user := domain.UserInfo{
		Email:    email,
		Username: name,
	}

	return user, nil
}
