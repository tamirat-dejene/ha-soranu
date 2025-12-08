package client

import (
	"fmt"

	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client authpb.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	client := authpb.NewAuthServiceClient(conn)
	return &AuthClient{
		Client: client,
	}, nil
}
