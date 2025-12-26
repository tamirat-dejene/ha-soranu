package valkey

import (
	"context"
	"fmt"
	"time"

	"github.com/tamirat-dejene/ha-soranu/shared/pkg/caching"
	"github.com/valkey-io/valkey-go"
)

type valkeyClient struct {
	client valkey.Client
}

// Close implements [caching.CacheClient].
func (v *valkeyClient) Close() error {
	v.client.Close()
	return nil
}

// Decrement implements [caching.CacheClient].
func (v *valkeyClient) Decrement(ctx context.Context, key string) (int64, error) {
	cmd := v.client.B().Decr().Key(key).Build()
	
	val, err := v.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s: %w", key, err)
	}
	return val, nil
}

// Delete implements [caching.CacheClient].
func (v *valkeyClient) Delete(ctx context.Context, key string) error {
	cmd := v.client.B().Del().Key(key).Build()
	
	if err := v.client.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists implements [caching.CacheClient].
func (v *valkeyClient) Exists(ctx context.Context, key string) (bool, error) {
	cmd := v.client.B().Exists().Key(key).Build()
	
	count, err := v.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

// Expire implements [caching.CacheClient].
func (v *valkeyClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	seconds := int64(expiration.Seconds())
	cmd := v.client.B().Expire().Key(key).Seconds(seconds).Build()

	if err := v.client.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}

// Get implements [caching.CacheClient].
func (v *valkeyClient) Get(ctx context.Context, key string) (string, error) {
	cmd := v.client.B().Get().Key(key).Build()
	
	val, err := v.client.Do(ctx, cmd).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", fmt.Errorf("key %s does not exist", key)
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// GetClient implements [caching.CacheClient].
func (v *valkeyClient) GetClient() any {
	return v.client
}

// Increment implements [caching.CacheClient].
func (v *valkeyClient) Increment(ctx context.Context, key string) (int64, error) {
	cmd := v.client.B().Incr().Key(key).Build()
	
	val, err := v.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return val, nil
}

// Set implements [caching.CacheClient].
func (v *valkeyClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	strVal := fmt.Sprintf("%v", value)

	cmd := v.client.B().Set().Key(key).Value(strVal).Ex(expiration).Build()
	
	if err := v.client.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// NewValkeyClient initializes the Valkey client
func NewValkeyClient(host string, port int, password string, db int) (caching.CacheClient, error) {
	opts := valkey.ClientOption{
		InitAddress:  []string{fmt.Sprintf("%s:%d", host, port)},
		Password:     password,
		SelectDB:     db,
		DisableCache: true,
	}

	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	return &valkeyClient{
		client: client,
	}, nil
}