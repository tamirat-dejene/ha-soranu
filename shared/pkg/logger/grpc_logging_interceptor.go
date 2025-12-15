package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	duration := time.Since(start)
	code := status.Code(err)

	fields := []zap.Field{
		zap.String("method", info.FullMethod),
		zap.Duration("duration", duration),
		zap.String("code", code.String()),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		Error("gRPC Request Failed", fields...)
	} else {
		Info("gRPC Request Success", fields...)
	}

	return resp, err
}
