package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/migrations"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/caching/valkey"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"

	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v5/stdlib"
	authservice "github.com/tamirat-dejene/ha-soranu/services/auth-service"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/handler"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/repository"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/usecase"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
)

func main() {
	// 1. Load Configuration
	env, err := authservice.GetEnv()
	if err != nil {
		panic(err)
	}

	// 2. Initialize Logger
	logger.InitLogger(env.SRV_ENV)
	defer logger.Log.Sync()
	logger.Info("Auth Service is running...", zap.String("env", env.SRV_ENV))

	// 3. Initialize Migrations
	migrate := migrations.NewMigrator(*env)
	if err := migrate.Migrate(context.Background(), "/app/auth-svc-migrations"); err != nil {
		logger.Fatal("migrations failed", zap.Error(err))
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)

	// 4. Initialize Postgres Client
	pgClient, err := postgres.NewPostgresClient(dsn)
	if err != nil {
		logger.Fatal("Failed to connect to Postgres", zap.Error(err))
	}
	defer pgClient.Close()

	// 5. Initialize Valkey Redis Client
	valkeyClient, err := valkey.NewValkeyClient(env.RedisHOST, env.RedisPort, env.RedisPassword, env.RedisDB)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer valkeyClient.Close()
	// 6. Initialize Repositories
	userRepo := repository.NewUserRepository(pgClient)
	authRepo := repository.NewAuthRepository(valkeyClient, time.Duration(5)*time.Minute)

	// 7. Initialize Usecases
	timeout := time.Duration(5) * time.Second // Default timeout
	userUsecase := usecase.NewUserUsecase(userRepo, timeout)
	authUsecase := usecase.NewAuthUsecase(timeout, authRepo, userRepo, *env)

	// 8. Initialize gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.AUTH_SRV_PORT))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggingInterceptor),
	)

	// 9. Register Handlers
	handler.NewGrpcAuthHandler(s, authUsecase)
	handler.NewGrpcUserHandler(s, userUsecase)

	// 10. Start Server
	logger.Info("Auth Service listening", zap.String("port", env.AUTH_SRV_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
