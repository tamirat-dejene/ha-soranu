package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	internalutil "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/util"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/caching"
)

type authRepository struct {
	client     caching.CacheClient
	expiration time.Duration
}

func getRefreshKey(tokenID string) string {
	return fmt.Sprintf("refresh:%s", internalutil.HashToken(tokenID))
}

type RefreshMeta struct {
	Email string `json:"email"`
	// CreatedAt int64  `json:"created_at"`
	// ExpiresAt int64  `json:"expires_at"`
	// DeviceID  string `json:"device_id,omitempty"`
}

func (a *authRepository) SaveRefreshToken(ctx context.Context, email, tokenID string) error {
	key := getRefreshKey(tokenID)
	meta := RefreshMeta{Email: email}

	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	return a.client.Set(ctx, key, string(data), a.expiration)
}

func (a *authRepository) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	key := getRefreshKey(tokenID)
	return a.client.Delete(ctx, key)
}

// ValidateRefreshToken retrieves the user email associated with a token without revoking it
func (a *authRepository) ValidateRefreshToken(ctx context.Context, tokenID string) (string, error) {
	key := getRefreshKey(tokenID)

	exists, err := a.client.Exists(ctx, key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errs.ErrTokenRevoked
	}

	data, err := a.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	var meta RefreshMeta
	if err := json.Unmarshal([]byte(data), &meta); err != nil {
		return "", err
	}

	return meta.Email, nil
}

// ConsumeRefreshToken retrieves and deletes the refresh token (rotation-safe, one-time-use)
func (a *authRepository) ConsumeRefreshToken(ctx context.Context, tokenId string) (string, error) {
	key := getRefreshKey(tokenId)

	// Check if key exists first
	exists, err := a.client.Exists(ctx, key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errs.ErrTokenRevoked
	}

	data, err := a.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	// Delete the token immediately for one-time use
	if err := a.client.Delete(ctx, key); err != nil {
		return "", err
	}

	var meta RefreshMeta
	if err := json.Unmarshal([]byte(data), &meta); err != nil {
		return "", err
	}

	return meta.Email, nil
}

// NewAuthRepository creates a Redis-based AuthRepository
func NewAuthRepository(client caching.CacheClient, expiration time.Duration) domain.AuthRepository {
	return &authRepository{
		client:     client,
		expiration: expiration,
	}
}
