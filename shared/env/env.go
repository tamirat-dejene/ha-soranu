package env

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	// Service settings
	SRV_ENV       string `mapstructure:"SRV_ENV"`
	AUTH_SRV_NAME string `mapstructure:"AUTH_SRV_NAME"`

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
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// Other environment variables can be added here
}

func LoadEnv() (*Env, error) {
	v := viper.New()
	v.AutomaticEnv()

	var env Env
	if err := v.Unmarshal(&env); err != nil {
		return nil, err
	}

	log.Printf("Service running in %s mode", env.SRV_ENV)
	return &env, nil
}