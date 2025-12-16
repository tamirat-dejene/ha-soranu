package usecase

import (
	"context"
	"time"

	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
)

type restaurantUseCase struct {
	repo    domain.RestaurantRepository
	timeout time.Duration
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

func NewRestaurantUseCase(repo domain.RestaurantRepository, timeout time.Duration) domain.RestaurantUseCase {
	return &restaurantUseCase{repo: repo, timeout: timeout}
}
