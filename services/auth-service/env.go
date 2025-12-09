package authservice

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	// Service settings
	SRV_ENV       string `mapstructure:"SRV_ENV"`
	AUTH_SRV_NAME string `mapstructure:"AUTH_SRV_NAME"`
	AUTH_SRV_PORT string `mapstructure:"AUTH_SRV_PORT"`

	// Google OAuth2 settings
	GoogleClientID string `mapstructure:"GOOGLE_CLIENT_ID"`

	// JWT settings
	AccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenTTL     string `mapstructure:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL    string `mapstructure:"REFRESH_TOKEN_TTL"`

	// Database settings
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`

	// Redis settings
	RedisHOST     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// Other environment variables can be added here
}

func getString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %s: %s, using default %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

func GetEnv() (*Env, error) {
	env := Env{
		SRV_ENV:            getString("SRV_ENV", "development"),
		AUTH_SRV_NAME:      getString("AUTH_SRV_NAME", "auth-service"),
		AUTH_SRV_PORT:      getString("AUTH_SRV_PORT", "9090"),
		GoogleClientID:     getString("GOOGLE_CLIENT_ID", ""),
		AccessTokenSecret:  getString("ACCESS_TOKEN_SECRET", "default_access_secret"),
		RefreshTokenSecret: getString("REFRESH_TOKEN_SECRET", "default_refresh_secret"),
		AccessTokenTTL:     getString("ACCESS_TOKEN_TTL", "15m"),
		RefreshTokenTTL:    getString("REFRESH_TOKEN_TTL", "7d"),
		DBHost:             getString("POSTGRES_HOST", "localhost"),
		DBPort:             getString("POSTGRES_PORT", "5432"),
		DBUser:             getString("POSTGRES_USER", "postgres"),
		DBPassword:         getString("POSTGRES_PASSWORD", "password"),
		DBName:             getString("POSTGRES_DB", "authdb"),
		RedisHOST:          getString("REDIS_HOST", "localhost"),
		RedisPort:          getInt("REDIS_PORT", 6379),
		RedisPassword:      getString("REDIS_PASSWORD", ""),
		RedisDB:            getInt("REDIS_DB", 0),
	}
	return &env, nil
}
