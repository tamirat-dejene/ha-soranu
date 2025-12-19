package usecase

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
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
	timeout  time.Duration
}

// PlaceOrder implements [domain.RestaurantUseCase].
func (r *restaurantUseCase) PlaceOrder(ctx context.Context, order *domain.PlaceOrder) (*domain.Order, error) {
	c, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	orderId := bytes.Buffer{}
	orderId.WriteString(order.CustomerID)
	orderId.WriteString("-")
	orderId.WriteString(order.RestaurantID)
	orderId.WriteString("-")
	orderId.WriteString(time.Now().Format("20060102150405"))

	// Make order processing logic here: verify items, calculate total, etc.
	// If we save order to DB, we can get the order ID from there.

	// Create domain event
	create_event := orderpb.OrderCreated{
		OrderId:       orderId.String(),
		CustomerId:    order.CustomerID,
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
		Key:   orderId.Bytes(),
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

	return &domain.Order{
		OrderID: orderId.String(),
		Status:  domain.ORDER_STATUS_PENDING,
	}, nil
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
	producer kafka.Producer, timeout time.Duration) domain.RestaurantUseCase {
	return &restaurantUseCase{repo: repo, producer: producer, timeout: timeout}
}
