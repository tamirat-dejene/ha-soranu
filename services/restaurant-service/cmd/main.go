package main

import (
	"context"
	"fmt"
	"net"

	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	svc "github.com/tamirat-dejene/ha-soranu/services/restaurant-service"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/api/grpc/handler"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/repository"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/usecase"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/migrations"
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
	logger.Info("restaurant-service service is starting...", zap.String("env", env.SRV_ENV))

	// 3. Migrate Database
	migrator := migrations.NewMigrator(*env)

	if err := migrator.Migrate(context.Background(), "/app/restaurant-svc-migrations"); err != nil {
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

	// 5. Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.RESTAURANT_SRV_PORT))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	
	// 6. Initialize Repository, Usecase, and register Handler
	restaurant_repo := repository.NewRestaurantRepository(pgClient)
	restaurant_usecase := usecase.NewRestaurantUseCase(restaurant_repo, 5)
	handler.NewRestaurantHandler(s, restaurant_usecase)

	logger.Info("Service listening", zap.String("port", env.RESTAURANT_SRV_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
