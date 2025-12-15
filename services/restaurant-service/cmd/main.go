package main

import (
	"fmt"
	"net"

	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	svc "github.com/tamirat-dejene/ha-soranu/services/restaurant-service"
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

	// 3. Start gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.RESTAURANT_SRV_PORT))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	// TODO: Register protobuf services

	logger.Info("Service listening", zap.String("port", env.RESTAURANT_SRV_PORT))
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
