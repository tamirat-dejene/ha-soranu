package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	authservice "github.com/tamirat-dejene/ha-soranu/services/auth-service"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/handler"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/repository"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/usecase"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
	"github.com/tamirat-dejene/ha-soranu/shared/redis"
)

func main() {
	fmt.Println("Auth Service is running...")

	env, err := authservice.GetEnv()
	if err != nil {
		panic(err)
	}

	// 1. Initialize Postgres
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	pgClient, err := postgres.NewPostgresClient(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgClient.Close()

	// 2. Initialize Redis
	redisClient := redis.NewRedisClient(env.RedisHOST, env.RedisPort, env.RedisPassword, env.RedisDB)
	defer redisClient.Close()

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(pgClient)

	// 4. Initialize Usecases
	timeout := time.Duration(5) * time.Second // Default timeout
	userUsecase := usecase.NewUserUsecase(userRepo, timeout)
	authUsecase := usecase.NewAuthUsecase(userUsecase, redisClient, timeout, *env)

	// 5. Initialize gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.AUTH_SRV_PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// 6. Register Handlers
	handler.NewGrpcAuthHandler(s, authUsecase)
	handler.NewGrpcUserHandler(s, userUsecase)

	// 7. Start Server
	log.Printf("Auth Service listening on port %s", env.AUTH_SRV_PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
