package notificationservice

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	// Service settings
	SRV_ENV                         string `mapstructure:"SRV_ENV"`
	NOTIFICATION_SRV_PORT           string `mapstructure:"NOTIFICATION_SRV_PORT"`
	NOTIFICATION_SRV_CONSUMER_GROUP string `mapstructure:"NOTIFICATION_SRV_CONSUMER_GROUP"`

	// Database settings
	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`

	// Kafka settings
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
		SRV_ENV:                         getString("SRV_ENV", "development"),
		NOTIFICATION_SRV_PORT:           getString("NOTIFICATION_SRV_PORT", "50053"),
		NOTIFICATION_SRV_CONSUMER_GROUP: getString("NOTIFICATION_SRV_CONSUMER_GROUP", "notification-service-group"),
		DBHost:                          getString("POSTGRES_HOST", "postgres-db"),
		DBPort:                          getString("POSTGRES_PORT", "5432"),
		DBUser:                          getString("POSTGRES_USER", "postgres"),
		DBPassword:                      getString("POSTGRES_PASSWORD", "password"),
		DBName:                          getString("POSTGRES_DB", "notification_db"),
		KafkaBroker:                     getString("KAFKA_BROKER_URL", "localhost:9092"),
	}
	return &env, nil
}
