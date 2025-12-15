package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

// LoggingInterceptor is a placeholder unary interceptor.
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// TODO: add request logging here
	return handler(ctx, req)
}
