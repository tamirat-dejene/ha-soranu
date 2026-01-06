package paymentservice

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	// Service settings
	SRV_ENV           string `mapstructure:"SRV_ENV"`
	PAYMENT_SRV_NAME  string `mapstructure:"PAYMENT_SRV_NAME"`
	PAYMENT_SRV_PORT  string `mapstructure:"PAYMENT_SRV_PORT"`
	PAYMENT_HTTP_PORT string `mapstructure:"PAYMENT_HTTP_PORT"`

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

	// Stripe settings
	StripeSecretKey     string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET"`

	// Kafka settings (for publishing order status updates)
	KafkaBroker string `mapstructure:"KAFKA_BROKER_URL"`
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
		SRV_ENV:             getString("SRV_ENV", "development"),
		PAYMENT_SRV_NAME:    getString("PAYMENT_SRV_NAME", "payment-service"),
		PAYMENT_SRV_PORT:    getString("PAYMENT_SRV_PORT", "9090"),
		PAYMENT_HTTP_PORT:   getString("PAYMENT_HTTP_PORT", "8081"),
		DBHost:              getString("POSTGRES_HOST", "postgres-db"),
		DBPort:              getString("POSTGRES_PORT", "5432"),
		DBUser:              getString("POSTGRES_USER", "postgres"),
		DBPassword:          getString("POSTGRES_PASSWORD", "password"),
		DBName:              getString("POSTGRES_DB", "payment-servicedb"),
		RedisHOST:           getString("REDIS_HOST", "localhost"),
		RedisPort:           getInt("REDIS_PORT", 6379),
		RedisPassword:       getString("REDIS_PASSWORD", ""),
		RedisDB:             getInt("REDIS_DB", 0),
		StripeSecretKey:     getString("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getString("STRIPE_WEBHOOK_SECRET", ""),
		KafkaBroker:         getString("KAFKA_BROKER_URL", "localhost:9092"),
	}
	return &env, nil
}
