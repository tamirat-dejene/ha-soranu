package apigateway

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	// Service settings
	SRV_ENV          string `mapstructure:"SRV_ENV"`
	API_GATEWAY_PORT string `mapstructure:"API_GATEWAY_PORT"`

	// Auth Service settings
	AUTH_SRV_NAME string `mapstructure:"AUTH_SRV_NAME"`
	AUTH_SRV_PORT string `mapstructure:"AUTH_SRV_PORT"`

	// Restaurant Service settings
	RESTAURANT_SRV_NAME string `mapstructure:"RESTAURANT_SRV_NAME"`
	RESTAURANT_SRV_PORT string `mapstructure:"RESTAURANT_SRV_PORT"`

	// Notification Service settings
	NOTIFICATION_SRV_NAME string `mapstructure:"NOTIFICATION_SRV_NAME"`
	NOTIFICATION_SRV_PORT string `mapstructure:"NOTIFICATION_SRV_PORT"`
}

func getString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func _(key string, defaultValue int) int {
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
		SRV_ENV:          getString("SRV_ENV", "development"),
		AUTH_SRV_NAME:    getString("AUTH_SRV_NAME", "auth-service"),
		AUTH_SRV_PORT:    getString("AUTH_SRV_PORT", "9090"),
		API_GATEWAY_PORT: getString("API_GATEWAY_PORT", "8080"),

		RESTAURANT_SRV_NAME: getString("RESTAURANT_SRV_NAME", "restaurant-service"),
		RESTAURANT_SRV_PORT: getString("RESTAURANT_SRV_PORT", "9091"),

		NOTIFICATION_SRV_NAME: getString("NOTIFICATION_SRV_NAME", "notification-service"),
		NOTIFICATION_SRV_PORT: getString("NOTIFICATION_SRV_PORT", "50053"),
	}

	return &env, nil
}
