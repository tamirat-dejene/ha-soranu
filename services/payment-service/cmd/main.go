package main

import (
	"context"
	"fmt"
	"net"
	nethttp "net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stripe/stripe-go/v78"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	svc "github.com/tamirat-dejene/ha-soranu/services/payment-service"
	httpapi "github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/api/http"
	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/repository"
	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/usecase"
	"github.com/tamirat-dejene/ha-soranu/services/payment-service/migrations"
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
	logger.Info("payment-service is starting...", zap.String("env", env.SRV_ENV))

	// 3. Run DB migrations
	migrator := migrations.NewMigrator(*env)
	if err := migrator.Migrate(context.Background(), "/app/payment-svc-migrations"); err != nil {
		logger.Fatal("failed to migrate database", zap.Error(err))
	}

	// 4. Initialize Postgres client
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", env.DBUser, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	pgClient, err := postgres.NewPostgresClient(dsn)
	if err != nil {
		logger.Fatal("failed to connect to Postgres", zap.Error(err))
	}
	defer pgClient.Close()

	// 5. Stripe API key
	if env.StripeSecretKey == "" {
		logger.Warn("STRIPE_SECRET_KEY not set; payment intent creation will fail")
	}
	stripe.Key = env.StripeSecretKey

	// 6. Wire repository and usecase
	repo := repository.NewPostgresRepository(pgClient)
	uc, err := usecase.New(repo, 10*time.Second, env.KafkaBroker)
	if err != nil {
		logger.Fatal("failed to init usecase", zap.Error(err))
	}

	// 7. Start HTTP server for intents + webhooks
	httpSrv := httpapi.NewServer(uc, env.StripeWebhookSecret)
	go func() {
		addr := ":" + env.PAYMENT_HTTP_PORT
		logger.Info("Payment HTTP server listening", zap.String("port", env.PAYMENT_HTTP_PORT))
		if err := nethttp.ListenAndServe(addr, httpSrv.Routes()); err != nil {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()

	// 8. Start gRPC server (reserved for future payment RPCs)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.PAYMENT_SRV_PORT))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	logger.Info("Payment gRPC server listening", zap.String("port", env.PAYMENT_SRV_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
