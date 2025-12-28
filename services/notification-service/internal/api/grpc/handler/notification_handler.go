package handler

import (
	"context"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/notification-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	notificationpb "github.com/tamirat-dejene/ha-soranu/shared/protos/notificationpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type notificationHandler struct {
	notificationpb.UnimplementedNotificationServiceServer
	usecase domain.NotificationUseCase
}

func NewNotificationHandler(server *grpc.Server, usecase domain.NotificationUseCase) {
	handler := &notificationHandler{
		usecase: usecase,
	}
	notificationpb.RegisterNotificationServiceServer(server, handler)
}

func (h *notificationHandler) GetNotifications(ctx context.Context, req *notificationpb.GetNotificationsRequest) (*notificationpb.GetNotificationsResponse, error) {
	if req == nil || req.RecipientId == "" || req.RecipientType == "" {
		logger.Error("invalid request")
		return nil, domain.ErrInvalidRequest
	}

	notifications, err := h.usecase.GetNotifications(ctx, req.RecipientId, req.RecipientType)
	if err != nil {
		logger.Error("failed to get notifications", zap.Error(err))
		return nil, err
	}

	var notificationProtos []*notificationpb.Notification
	for _, n := range notifications {
		notificationProtos = append(notificationProtos, &notificationpb.Notification{
			Id:            n.ID,
			RecipientId:   n.RecipientID,
			RecipientType: n.RecipientType,
			OrderId:       n.OrderID,
			Title:         n.Title,
			Message:       n.Message,
			IsRead:        n.IsRead,
			Type:          n.Type,
			CreatedAt:     n.CreatedAt.Format(time.RFC3339),
		})
	}

	return &notificationpb.GetNotificationsResponse{
		Notifications: notificationProtos,
	}, nil
}

func (h *notificationHandler) MarkAsRead(ctx context.Context, req *notificationpb.MarkAsReadRequest) (*notificationpb.MarkAsReadResponse, error) {
	if req == nil || req.NotificationId == "" {
		logger.Error("invalid request")
		return nil, domain.ErrInvalidRequest
	}

	err := h.usecase.MarkAsRead(ctx, req.NotificationId)
	if err != nil {
		logger.Error("failed to mark notification as read", zap.Error(err))
		return &notificationpb.MarkAsReadResponse{Success: false}, err
	}

	return &notificationpb.MarkAsReadResponse{Success: true}, nil
}
