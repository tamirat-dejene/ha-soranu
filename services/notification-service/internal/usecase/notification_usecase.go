package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/notification-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/events"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka"
	envent_envelope "github.com/tamirat-dejene/ha-soranu/shared/protos/envent_envelopepb"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/orderpb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type notificationUseCase struct {
	repo     domain.NotificationRepository
	consumer kafka.Consumer
	timeout  time.Duration
}

func NewNotificationUseCase(
	repo domain.NotificationRepository,
	consumer kafka.Consumer,
	timeout time.Duration,
) domain.NotificationUseCase {
	return &notificationUseCase{
		repo:     repo,
		consumer: consumer,
		timeout:  timeout,
	}
}

func (uc *notificationUseCase) StartConsumer(ctx context.Context) error {
	logger.Info("Starting notification service Kafka consumer")

	return uc.consumer.Subscribe(ctx, []string{
		events.OrderPlacedEvent,
		events.OrderStatusUpdatedEvent,
	}, func(msgCtx context.Context, msg *kafka.Message) error {
		// 1. Unmarshal envelope using binary protobuf
		var envelope envent_envelope.EventEnvelope
		if err := proto.Unmarshal(msg.Value, &envelope); err != nil {
			logger.Error("failed to unmarshal event envelope", zap.Error(err))
			return err
		}

		logger.Info("Received event",
			zap.String("event_type", envelope.EventType),
			zap.String("event_id", envelope.EventId),
			zap.String("message_key", string(msg.Key)))

		// 2. Route by envelope.EventType
		switch envelope.EventType {
		case events.OrderPlacedEvent:
			return uc.handleOrderPlaced(msgCtx, &envelope)
		case events.OrderStatusUpdatedEvent:
			return uc.handleOrderStatusUpdated(msgCtx, &envelope)
		default:
			logger.Warn("unknown event type", zap.String("event_type", envelope.EventType))
			return nil
		}
	})
}

func (uc *notificationUseCase) handleOrderPlaced(ctx context.Context, envelope *envent_envelope.EventEnvelope) error {
	// Unmarshal payload using binary protobuf
	var orderCreated orderpb.OrderCreated
	if err := proto.Unmarshal(envelope.Payload, &orderCreated); err != nil {
		logger.Error("failed to unmarshal OrderCreated event", zap.Error(err))
		return err
	}

	logger.Info("Processing OrderPlaced event",
		zap.String("order_id", orderCreated.OrderId),
		zap.String("customer_id", orderCreated.CustomerId),
		zap.String("restaurant_id", orderCreated.RestaurantId))

	// Create notification for restaurant
	notification := &domain.Notification{
		RecipientID:   orderCreated.RestaurantId,
		RecipientType: "RESTAURANT",
		OrderID:       orderCreated.OrderId,
		Title:         "New Order Received",
		Message:       fmt.Sprintf("You have a new order #%s for $%.2f", orderCreated.OrderId[:8], orderCreated.TotalAmount),
		IsRead:        false,
		Type:          "NEW_ORDER",
	}

	if err := uc.repo.CreateNotification(ctx, notification); err != nil {
		logger.Error("failed to create notification for restaurant", zap.Error(err))
		return err
	}

	return nil
}

func (uc *notificationUseCase) handleOrderStatusUpdated(ctx context.Context, envelope *envent_envelope.EventEnvelope) error {
	// Unmarshal payload using binary protobuf
	var orderStatusUpdated orderpb.OrderStatusUpdated
	if err := proto.Unmarshal(envelope.Payload, &orderStatusUpdated); err != nil {
		logger.Error("failed to unmarshal OrderStatusUpdated event", zap.Error(err))
		return err
	}

	logger.Info("Processing OrderStatusUpdated event",
		zap.String("order_id", orderStatusUpdated.OrderId),
		zap.String("customer_id", orderStatusUpdated.CustomerId),
		zap.String("new_status", orderStatusUpdated.NewStatus.String()))

	// Create notification for customer
	notification := &domain.Notification{
		RecipientID:   orderStatusUpdated.CustomerId,
		RecipientType: "USER",
		OrderID:       orderStatusUpdated.OrderId,
		Title:         "Order Status Updated",
		Message:       fmt.Sprintf("Your order status has been updated to: %s", orderStatusUpdated.NewStatus.String()),
		IsRead:        false,
		Type:          "ORDER_UPDATE",
	}

	if err := uc.repo.CreateNotification(ctx, notification); err != nil {
		logger.Error("failed to create notification for user", zap.Error(err))
		return err
	}

	return nil
}

func (uc *notificationUseCase) GetNotifications(ctx context.Context, recipientID string, recipientType string) ([]*domain.Notification, error) {
	c, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.GetNotifications(c, recipientID, recipientType)
}

func (uc *notificationUseCase) MarkAsRead(ctx context.Context, notificationID string) error {
	c, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.MarkAsRead(c, notificationID)
}

func (uc *notificationUseCase) DeleteNotification(ctx context.Context, notificationID string) error {
	c, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	return uc.repo.DeleteNotification(c, notificationID)
}