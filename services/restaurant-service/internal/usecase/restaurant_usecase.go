package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/events"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka"
	envent_envelope "github.com/tamirat-dejene/ha-soranu/shared/protos/envent_envelopepb"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/orderpb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type restaurantUseCase struct {
	repo     domain.RestaurantRepository
	producer kafka.Producer
	consumer kafka.Consumer
	timeout  time.Duration
}

// ShipOrder implements [domain.RestaurantUseCase].
func (r *restaurantUseCase) ShipOrder(ctx context.Context, restaurantID string, orderID string) (string, string, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.ShipOrder(c, restaurantID, orderID)
}

// UpdateOrderStatus implements [domain.RestaurantUseCase].
func (r *restaurantUseCase) UpdateOrderStatus(ctx context.Context, restaurantID string, orderID string, newStatus string) (*domain.Order, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	ord, err := r.repo.UpdateOrderStatus(c, restaurantID, orderID, newStatus)
	if err != nil {
		return nil, err
	}

	// Create domain event for order status update
	update_event := orderpb.OrderStatusUpdated{
		OrderId:       ord.OrderId,
		NewStatus:     dto.DomainOrderStatusToProto(ord.Status),
		UpdatedAtUnix: time.Now().Unix(),
	}

	// Serialize event
	eventData, err := protojson.Marshal(&update_event)
	if err != nil {
		logger.Error("failed to marshal order status updated event", zap.Error(err))
		return nil, err
	}

	eventEnvelop := &envent_envelope.EventEnvelope{
		EventId:        uuid.NewString(),
		OccurredAtUnix: time.Now().Unix(),
		Payload:        eventData,
	}
	envelopeBytes, err := protojson.Marshal(eventEnvelop)
	if err != nil {
		logger.Error("failed to marshal event envelope", zap.Error(err))
		return nil, err
	}

	// Publish event to Kafka
	err = r.producer.Publish(c, &kafka.Message{
		Topic: events.OrderStatusUpdatedEvent,
		Key:   []byte(ord.OrderId),
		Value: envelopeBytes,
		Headers: map[string][]byte{
			"event_type":   []byte(events.OrderStatusUpdatedEvent),
			"content_type": []byte("application/x-protobuf"),
			"producer":     []byte("restaurant-service"),
		},
	})

	if err != nil {
		logger.Error("failed to publish order status updated event", zap.Error(err))
		return nil, err
	}

	return ord, nil
}

// GetOrders implements [domain.RestaurantUseCase].
func (r *restaurantUseCase) GetOrders(ctx context.Context, restaurantID string) ([]domain.Order, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.GetOrders(c, restaurantID)
}

// PlaceOrder implements [domain.RestaurantUseCase].
func (r *restaurantUseCase) PlaceOrder(ctx context.Context, order *domain.PlaceOrder) (*domain.Order, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	ord, err := r.repo.PlaceOrder(c, order)
	if err != nil {
		return nil, err
	}

	// Create domain event
	create_event := orderpb.OrderCreated{
		OrderId:       ord.OrderId,
		CustomerId:    order.CustomerID,
		TotalAmount:   ord.TotalAmount,
		CreatedAtUnix: time.Now().Unix(),
	}

	// Serialize event
	eventData, err := protojson.Marshal(&create_event)
	if err != nil {
		logger.Error("failed to marshal order created event", zap.Error(err))
		return nil, err
	}

	eventEnvelop := &envent_envelope.EventEnvelope{
		EventId:        uuid.NewString(),
		OccurredAtUnix: time.Now().Unix(),
		Payload:        eventData,
	}
	envelopeBytes, err := protojson.Marshal(eventEnvelop)
	if err != nil {
		logger.Error("failed to marshal event envelope", zap.Error(err))
		return nil, err
	}

	// Publish event to Kafka

	err = r.producer.Publish(c, &kafka.Message{
		Topic: events.OrderPlacedEvent,
		Key:   []byte(ord.OrderId),
		Value: envelopeBytes,
		Headers: map[string][]byte{
			"event_type":   []byte(events.OrderPlacedEvent),
			"content_type": []byte("application/x-protobuf"),
			"producer":     []byte("restaurant-service"),
		},
	})

	if err != nil {
		logger.Log.Error("failed to publish order placed event", zap.Error(err))
		return nil, err
	}

	return ord, nil
}

// LoginRestaurant implements domain.RestaurantUseCase.
func (r *restaurantUseCase) LoginRestaurant(ctx context.Context, email string, secretKey string) (*domain.Restaurant, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.LoginRestaurant(c, email, secretKey)
}

// AddMenuItem implements domain.RestaurantUseCase.
func (r *restaurantUseCase) AddMenuItem(ctx context.Context, restaurantID string, item domain.MenuItem) (*domain.MenuItem, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.AddMenuItem(c, restaurantID, item)
}

// GetRestaurantByID implements domain.RestaurantUseCase.
func (r *restaurantUseCase) GetRestaurantByID(ctx context.Context, restaurantID string) (*domain.Restaurant, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.GetRestaurantByID(c, restaurantID)
}

// GetRestaurants implements domain.RestaurantUseCase.
func (r *restaurantUseCase) StreamRestaurants(
	ctx context.Context,
	area domain.Area,
	onResult func(domain.Restaurant) error,
) error {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	return r.repo.StreamRestaurants(c, area, onResult)
}

// RegisterRestaurant implements domain.RestaurantUseCase.
func (r *restaurantUseCase) RegisterRestaurant(ctx context.Context, restaurant *domain.Restaurant) (*domain.Restaurant, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.CreateRestaurant(c, restaurant)
}

// RemoveMenuItem implements domain.RestaurantUseCase.
func (r *restaurantUseCase) RemoveMenuItem(ctx context.Context, restaurantID string, itemID string) error {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.RemoveMenuItem(c, restaurantID, itemID)
}

// UpdateMenuItem implements domain.RestaurantUseCase.
func (r *restaurantUseCase) UpdateMenuItem(ctx context.Context, restaurantID string, item domain.MenuItem) (*domain.MenuItem, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.repo.UpdateMenuItem(c, restaurantID, item)
}

func NewRestaurantUseCase(repo domain.RestaurantRepository,
	producer kafka.Producer,
	consumer kafka.Consumer,
	timeout time.Duration) domain.RestaurantUseCase {
	return &restaurantUseCase{repo: repo, producer: producer, consumer: consumer, timeout: timeout}
}
