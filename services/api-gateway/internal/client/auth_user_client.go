package client

import (
	"fmt"

	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// User and Auth Service Client
type UAServiceClient struct {
	AuthClient authpb.AuthServiceClient
	UserClient userpb.UserServiceClient
	conn       *grpc.ClientConn
}

// New User and Auth Service Client
func NewUAServiceClient(addr string) (*UAServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	auth_client := authpb.NewAuthServiceClient(conn)
	user_client := userpb.NewUserServiceClient(conn)

	return &UAServiceClient{
		AuthClient: auth_client,
		UserClient: user_client,
	}, nil
}

// Close the gRPC connection
func (c *UAServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.Error("failed to close gRPC connection", zap.Error(err))
		}
	}
}
