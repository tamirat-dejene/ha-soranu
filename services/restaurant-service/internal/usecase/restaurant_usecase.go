package usecase

import (
	"context"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/events"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/orderpb"
	"go.uber.org/zap"
)

type restaurantUseCase struct {
	repo      domain.RestaurantRepository
	publisher events.EventPublisher
	timeout   time.Duration
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

	// Create and publish domain event
	update_event := orderpb.OrderStatusUpdated{
		OrderId:       ord.OrderId,
		CustomerId:    ord.CustomerID,
		NewStatus:     dto.DomainOrderStatusToProto(ord.Status),
		UpdatedAtUnix: time.Now().Unix(),
	}

	if err := r.publisher.PublishOrderStatusUpdated(c, &update_event); err != nil {
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

	// Create and publish domain event
	create_event := orderpb.OrderCreated{
		OrderId:       ord.OrderId,
		CustomerId:    order.CustomerID,
		RestaurantId:  ord.RestaurantID,
		TotalAmount:   ord.TotalAmount,
		CreatedAtUnix: time.Now().Unix(),
	}

	if err := r.publisher.PublishOrderCreated(c, &create_event); err != nil {
		logger.Error("failed to publish order created event", zap.Error(err))
		return nil, err
	}

	logger.Info("published order created event", zap.String("order_id", ord.OrderId), zap.String("restaurant_id", ord.RestaurantID), zap.Float64("total_amount", ord.TotalAmount))

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
	publisher events.EventPublisher,
	timeout time.Duration) domain.RestaurantUseCase {
	return &restaurantUseCase{repo: repo, publisher: publisher, timeout: timeout}
}
