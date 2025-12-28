package client

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	notificationpb "github.com/tamirat-dejene/ha-soranu/shared/protos/notificationpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationServiceClient struct {
	conn               *grpc.ClientConn
	NotificationClient notificationpb.NotificationServiceClient
}

func NewNotificationServiceClient(notificationServiceURL string) (*NotificationServiceClient, error) {
	conn, err := grpc.NewClient(notificationServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to Notification Service", zap.Error(err))
		return nil, err
	}

	client := notificationpb.NewNotificationServiceClient(conn)

	logger.Info("Connected to Notification Service", zap.String("url", notificationServiceURL))

	return &NotificationServiceClient{
		conn:               conn,
		NotificationClient: client,
	}, nil
}

func (c *NotificationServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.Error("failed to close gRPC connection", zap.Error(err))
		}
	}
}

func (c *NotificationServiceClient) GetNotifications(ctx context.Context, recipientID string, recipientType string) ([]*notificationpb.Notification, error) {
	req := &notificationpb.GetNotificationsRequest{
		RecipientId:   recipientID,
		RecipientType: recipientType,
	}
	resp, err := c.NotificationClient.GetNotifications(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Notifications, nil
}

func (c *NotificationServiceClient) MarkAsRead(ctx context.Context, notificationID string) (bool, error) {
	req := &notificationpb.MarkAsReadRequest{
		NotificationId: notificationID,
	}
	resp, err := c.NotificationClient.MarkAsRead(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
