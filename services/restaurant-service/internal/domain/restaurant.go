package domain

import (
	"context"
)

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

type Area struct {
	LatitudeMin float32
	LatitudeMax float32
	RadiusInKm  float32
}

type PlaceOrder struct {
	CustomerID   string
	RestaurantID string
	Items        []OrderItem
}

type Order struct {
	OrderId      string
	CustomerID   string
	RestaurantID string
	Items        []OrderItem
	TotalAmount  float64
	Status       string
}

type OrderItem struct {
	ItemId   string
	Quantity int32
}

type RestaurantUseCase interface {
	LoginRestaurant(ctx context.Context, email, secretKey string) (*Restaurant, error)

	RegisterRestaurant(ctx context.Context, restaurant *Restaurant) (*Restaurant, error)
	StreamRestaurants(ctx context.Context, area Area, onResult func(Restaurant) error) error
	GetRestaurantByID(ctx context.Context, restaurantID string) (*Restaurant, error)

	AddMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)
	RemoveMenuItem(ctx context.Context, restaurantID, itemID string) error
	UpdateMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)

	PlaceOrder(ctx context.Context, order *PlaceOrder) (*Order, error)
	GetOrders(ctx context.Context, restaurantID string) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, restaurantID, orderID, newStatus string) (*Order, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)

	ShipOrder(ctx context.Context, restaurantID, orderID string) (string, string, error)
}

type RestaurantRepository interface {
	LoginRestaurant(ctx context.Context, email, secretKey string) (*Restaurant, error)

	CreateRestaurant(ctx context.Context, restaurant *Restaurant) (*Restaurant, error)
	StreamRestaurants(ctx context.Context, area Area, onRow func(Restaurant) error) error
	GetRestaurantByID(ctx context.Context, restaurantID string) (*Restaurant, error)

	AddMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)
	RemoveMenuItem(ctx context.Context, restaurantID, itemID string) error
	UpdateMenuItem(ctx context.Context, restaurantID string, item MenuItem) (*MenuItem, error)

	PlaceOrder(ctx context.Context, order *PlaceOrder) (*Order, error)
	GetOrders(ctx context.Context, restaurantID string) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, restaurantID, orderID, newStatus string) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)

	GetOrder(ctx context.Context, orderID string) (*Order, error)
	ShipOrder(ctx context.Context, restaurantID, orderID string) (string, string, error)
}
