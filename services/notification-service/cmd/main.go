package main

import (
	"context"
	"fmt"
	"net"
	"time"

	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka/sarama"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	svc "github.com/tamirat-dejene/ha-soranu/services/notification-service"
	"github.com/tamirat-dejene/ha-soranu/services/notification-service/internal/api/grpc/handler"
	"github.com/tamirat-dejene/ha-soranu/services/notification-service/internal/repository"
	"github.com/tamirat-dejene/ha-soranu/services/notification-service/internal/usecase"
	"github.com/tamirat-dejene/ha-soranu/services/notification-service/migrations"
)

func main() {
	// 1. Load Configuration
	env, err := svc.GetEnv()
	if err != nil {
		panic(err)
	}

	// 2. Initialize Logger
	logger.InitLogger(env.SRV_ENV)
	defer logger.Log.Sync()
	logger.Info("notification-service is starting...", zap.String("env", env.SRV_ENV))

	// 3. Migrate Database
	migrator := migrations.NewMigrator(*env)

	if err := migrator.Migrate(context.Background(), "/app/notification-svc-migrations"); err != nil {
		logger.Fatal("failed to migrate database", zap.Error(err))
	}

	// 4. Initialize postgres client
	postgresDsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)

	pgClient, err := postgres.NewPostgresClient(postgresDsn)
	if err != nil {
		logger.Fatal("failed to connect to Postgres", zap.Error(err))
	}
	defer pgClient.Close()

	// 5. Initialize sarama consumer
	consumer, err := sarama.NewConsumer([]string{env.KafkaBroker}, env.NOTIFICATION_SRV_CONSUMER_GROUP)
	if err != nil {
		logger.Fatal("failed to create kafka consumer", zap.Error(err))
	}
	defer consumer.Close()

	// 6. Initialize Repository, Usecase, and register Handler
	notification_repo := repository.NewNotificationRepository(pgClient)
	notification_usecase := usecase.NewNotificationUseCase(notification_repo, consumer, 10*time.Second)

	// 7. Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.NOTIFICATION_SRV_PORT))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	s := grpc.NewServer()

	handler.NewNotificationHandler(s, notification_usecase)

	// Start Kafka Consumer for Notifications
	go func() {
		logger.Info("Starting notification consumer...")
		if err := notification_usecase.StartConsumer(context.Background()); err != nil {
			logger.Error("failed to start notification consumer", zap.Error(err))
		}
	}()

	logger.Info("Service listening", zap.String("port", env.NOTIFICATION_SRV_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
