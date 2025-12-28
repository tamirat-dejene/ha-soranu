package events

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka"
	envent_envelope "github.com/tamirat-dejene/ha-soranu/shared/protos/envent_envelopepb"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/orderpb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// EventPublisher handles publishing domain events to Kafka with proper protobuf serialization
type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, event *orderpb.OrderCreated) error
	PublishOrderStatusUpdated(ctx context.Context, event *orderpb.OrderStatusUpdated) error
}

type kafkaEventPublisher struct {
	producer kafka.Producer
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(producer kafka.Producer) EventPublisher {
	return &kafkaEventPublisher{
		producer: producer,
	}
}

// PublishOrderCreated publishes an OrderCreated event to Kafka
func (p *kafkaEventPublisher) PublishOrderCreated(ctx context.Context, event *orderpb.OrderCreated) error {
	return p.publishEvent(ctx, OrderPlacedEvent, event.OrderId, event)
}

// PublishOrderStatusUpdated publishes an OrderStatusUpdated event to Kafka
func (p *kafkaEventPublisher) PublishOrderStatusUpdated(ctx context.Context, event *orderpb.OrderStatusUpdated) error {
	return p.publishEvent(ctx, OrderStatusUpdatedEvent, event.OrderId, event)
}

// publishEvent is the generic event publishing logic
func (p *kafkaEventPublisher) publishEvent(ctx context.Context, eventType string, key string, event proto.Message) error {
	// 1. Marshal the domain event to binary protobuf
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		logger.Error("failed to marshal event", zap.String("event_type", eventType), zap.Error(err))
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 2. Create event envelope with event_type
	envelope := &envent_envelope.EventEnvelope{
		EventId:        uuid.NewString(),
		EventType:      eventType,
		OccurredAtUnix: time.Now().Unix(),
		Payload:        eventBytes,
	}

	// 3. Marshal the envelope to binary protobuf
	envelopeBytes, err := proto.Marshal(envelope)
	if err != nil {
		logger.Error("failed to marshal envelope", zap.String("event_type", eventType), zap.Error(err))
		return fmt.Errorf("failed to marshal envelope: %w", err)
	}

	// 4. Publish to Kafka with proper headers
	msg := &kafka.Message{
		Topic: eventType,
		Key:   []byte(key),
		Value: envelopeBytes,
		Headers: map[string][]byte{
			"event_type":   []byte(eventType),
			"content_type": []byte("application/x-protobuf"),
		},
	}

	if err := p.producer.Publish(ctx, msg); err != nil {
		logger.Error("failed to publish event to kafka",
			zap.String("event_type", eventType),
			zap.String("key", key),
			zap.Error(err))
		return fmt.Errorf("failed to publish event: %w", err)
	}

	logger.Info("published event to kafka",
		zap.String("event_type", eventType),
		zap.String("event_id", envelope.EventId),
		zap.String("key", key))

	return nil
}
