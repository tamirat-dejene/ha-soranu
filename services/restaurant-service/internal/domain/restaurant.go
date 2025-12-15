package domain

import "context"

type Restaurant struct {
	ID        string
	Email     string
	SecretKey string
	Name      string
	Latitude  float32
	Longitude float32
	MenuItems []MenuItem
}

type MenuItem struct {
	ItemID      string
	Name        string
	Description string
	Price       float32
}

type RestaurantUseCase interface {
	LoginRestaurant(ctx context.Context, email, secretKey string) (*Restaurant, error)

	RegisterRestaurant(ctx context.Context, restaurant *Restaurant) (*Restaurant, error)
	GetRestaurants(ctx context.Context) ([]Restaurant, error)
	GetRestaurantByID(ctx context.Context, restaurantID string) (*Restaurant, error)

	AddMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)
	RemoveMenuItem(ctx context.Context, restaurantID, itemID string) error
	UpdateMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)
}

type RestaurantRepository interface {
	LoginRestaurant(ctx context.Context, email, secretKey string) (*Restaurant, error)

	CreateRestaurant(ctx context.Context, restaurant *Restaurant) (*Restaurant, error)
	GetRestaurants(ctx context.Context) ([]Restaurant, error)
	GetRestaurantByID(ctx context.Context, restaurantID string) (*Restaurant, error)

	AddMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)
	RemoveMenuItem(ctx context.Context, restaurantID, itemID string) error
	UpdateMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)
}
