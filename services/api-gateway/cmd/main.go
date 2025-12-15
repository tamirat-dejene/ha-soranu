package main

import (
	apigateway "github.com/tamirat-dejene/ha-soranu/services/api-gateway"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/server"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 1. Load Configuration
	cfg, err := apigateway.GetEnv()
	if err != nil {
		// Fallback to standard logger if config fails
		panic("Failed to load config: " + err.Error())
	}

	// 2. Initialize Logger
	logger.InitLogger(cfg.SRV_ENV)
	defer logger.Log.Sync()
	logger.Info("Starting API Gateway...", zap.String("env", cfg.SRV_ENV))

	// 3. Initialize Auth Service Client
	uaServiceClient, err := client.NewUAServiceClient(cfg.AUTH_SRV_NAME + ":" + cfg.AUTH_SRV_PORT)
	if err != nil {
		logger.Fatal("Failed to connect to Auth Service", zap.Error(err))
	}
	logger.Info("Connected to Auth Service", zap.String("addr", cfg.AUTH_SRV_NAME+":"+cfg.AUTH_SRV_PORT))
	defer uaServiceClient.Close()

	// 4. Initialize Restaurant Service Client
	restaurantServiceClient, err := client.NewRestaurantServiceClient(cfg.RESTAURANT_SRV_NAME + ":" + cfg.RESTAURANT_SRV_PORT)
	if err != nil {
		logger.Fatal("Failed to connect to Restaurant Service", zap.Error(err))
	}
	logger.Info("Connected to Restaurant Service", zap.String("addr", cfg.RESTAURANT_SRV_NAME+":"+cfg.RESTAURANT_SRV_PORT))
	defer restaurantServiceClient.Close()

	// 4. Initialize and Run Server
	srv := server.NewServer(cfg, uaServiceClient, restaurantServiceClient)
	srv.SetupRoutes()

	logger.Info("API Gateway listening", zap.String("port", cfg.API_GATEWAY_PORT))
	if err := srv.Run(); err != nil {
		logger.Fatal("Server failed to run", zap.Error(err))
	}
}
