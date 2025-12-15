package client

import (
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/restaurantpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RestaurantServiceClient struct {
	RestaurantClient restaurantpb.RestaurantServiceClient
	conn *grpc.ClientConn
}

func NewRestaurantServiceClient(addr string) (*RestaurantServiceClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to gRPC server", zap.Error(err))
		return nil, err
	}

	client := restaurantpb.NewRestaurantServiceClient(conn)
	return &RestaurantServiceClient{
		RestaurantClient: client,
		conn:             conn,
	}, nil
}


func (c *RestaurantServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.Error("failed to close gRPC connection", zap.Error(err))
		}
	}
}