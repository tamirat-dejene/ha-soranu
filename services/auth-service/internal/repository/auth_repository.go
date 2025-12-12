package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	internalutil "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/util"
	"github.com/tamirat-dejene/ha-soranu/shared/redis"
)

type authRepository struct {
	client     redis.RedisClient
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

func (a *authRepository) SaveRefreshToken(email, tokenID string) error {
	key := getRefreshKey(tokenID)
	meta := RefreshMeta{Email: email}

	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	return a.client.Set(key, string(data), a.expiration)
}

func (a *authRepository) DeleteRefreshToken(tokenID string) error {
	key := getRefreshKey(tokenID)
	return a.client.Delete(key)
}

// ValidateRefreshToken retrieves the user email associated with a token without revoking it
func (a *authRepository) ValidateRefreshToken(tokenID string) (string, error) {
	key := getRefreshKey(tokenID)

	exists, err := a.client.Exists(key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errs.ErrTokenRevoked
	}

	data, err := a.client.Get(key)
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
func (a *authRepository) ConsumeRefreshToken(tokenId string) (string, error) {
	key := getRefreshKey(tokenId)

	// Check if key exists first
	exists, err := a.client.Exists(key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errs.ErrTokenRevoked
	}

	data, err := a.client.Get(key)
	if err != nil {
		return "", err
	}

	// Delete the token immediately for one-time use
	if err := a.client.Delete(key); err != nil {
		return "", err
	}

	var meta RefreshMeta
	if err := json.Unmarshal([]byte(data), &meta); err != nil {
		return "", err
	}

	return meta.Email, nil
}

// NewAuthRepository creates a Redis-based AuthRepository
func NewAuthRepository(client redis.RedisClient, expiration time.Duration) domain.AuthRepository {
	return &authRepository{
		client:     client,
		expiration: expiration,
	}
}
