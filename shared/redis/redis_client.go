package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient interface {
	GetClient() *redis.Client
	Close() error
	Set(key string, value any, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Increment(key string) (int64, error)
	Decrement(key string) (int64, error)
	Expire(key string, expiration time.Duration) error
}

type redisClient struct {
	client *redis.Client
}

type RedisService struct{}

func NewRedisClient(host string, port int, password string, db int) *redisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	return &redisClient{
		client: client,
	}
}

func (r *redisClient) GetClient() *redis.Client {
	return r.client
}

func (r *redisClient) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *redisClient) Set(key string, value any, expiration time.Duration) error {
	if err := r.client.Set(key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (r *redisClient) Get(key string) (string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s does not exist", key)
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return value, nil
}

func (r *redisClient) Delete(key string) error {
	if err := r.client.Del(key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func (r *redisClient) Exists(key string) (bool, error) {
	exists, err := r.client.Exists(key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return exists > 0, nil
}

func (r *redisClient) Increment(key string) (int64, error) {
	newValue, err := r.client.Incr(key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return newValue, nil
}

func (r *redisClient) Decrement(key string) (int64, error) {
	newValue, err := r.client.Decr(key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s: %w", key, err)
	}
	return newValue, nil
}

func (r *redisClient) Expire(key string, expiration time.Duration) error {
	if err := r.client.Expire(key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}
	return nil
}
