package main

import (
	apigateway "github.com/tamirat-dejene/ha-soranu/services/api-gateway"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/server"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	defer logger.Log.Sync()
	logger.Log.Info("Starting API Gateway...")

	// 2. Load Configuration
	cfg, err := apigateway.GetEnv()
	if err != nil {
		logger.Log.Fatal("Failed to load config", zap.Error(err))
	}

	// 3. Initialize Auth Service Client
	authClient, err := client.NewAuthClient(cfg.AUTH_SRV_NAME + ":" + cfg.AUTH_SRV_PORT)
	if err != nil {
		logger.Log.Fatal("Failed to connect to Auth Service", zap.Error(err))
	}
	logger.Log.Info("Connected to Auth Service", zap.String("addr", cfg.AUTH_SRV_NAME+":"+cfg.AUTH_SRV_PORT))

	// 4. Initialize and Run Server
	srv := server.NewServer(cfg, authClient)
	srv.SetupRoutes()

	logger.Log.Info("API Gateway listening", zap.String("port", cfg.API_GATEWAY_PORT))
	if err := srv.Run(); err != nil {
		logger.Log.Fatal("Server failed to run", zap.Error(err))
	}
}
