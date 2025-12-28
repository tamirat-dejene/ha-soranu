package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/restaurantpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type restaurantHandler struct {
	restaurantpb.UnimplementedRestaurantServiceServer
	restaurantUsecase domain.RestaurantUseCase
}

// ShipOrder implements [restaurantpb.RestaurantServiceServer].
func (r *restaurantHandler) ShipOrder(ctx context.Context, req *restaurantpb.ShipOrderRequest) (*restaurantpb.ShipOrderResponse, error) {
	if req == nil {
		return nil, domain.ErrInvalidOrderData
	}

	confirmation, driverID, err := r.restaurantUsecase.ShipOrder(ctx, req.RestaurantId, req.OrderId)
	if err != nil {
		return nil, err
	}

	logger.Info("shipped order", zap.String("order_id", req.OrderId), zap.String("driver_id", driverID))

	return &restaurantpb.ShipOrderResponse{
		ConfirmationMessage: confirmation,
		DriverId:            driverID,
	}, nil
}

// GetOrders implements [restaurantpb.RestaurantServiceServer].
func (r *restaurantHandler) GetOrders(ctx context.Context, req *restaurantpb.GetOrdersRequest) (*restaurantpb.GetOrdersResponse, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	orders, err := r.restaurantUsecase.GetOrders(ctx, req.RestaurantId)
	if err != nil {
		return nil, err
	}

	logger.Info("fetched orders", zap.Int("count", len(orders)))

	var orderProtos []*restaurantpb.Order
	for _, order := range orders {
		orderProtos = append(orderProtos, dto.DomainOrderToProto(order))
	}

	return &restaurantpb.GetOrdersResponse{
		Orders: orderProtos,
	}, nil
}

// UpdateOrderStatus implements [restaurantpb.RestaurantServiceServer].
func (r *restaurantHandler) UpdateOrderStatus(ctx context.Context, req *restaurantpb.UpdateOrderStatusRequest) (*restaurantpb.Order, error) {
	if req == nil {
		return nil, domain.ErrInvalidOrderData
	}

	updatedOrder, err := r.restaurantUsecase.UpdateOrderStatus(ctx, req.RestaurantId, req.OrderId, dto.ProtoOrderStatusToDomain(req.NewStatus))
	if err != nil {
		return nil, err
	}

	logger.Info("updated order status", zap.String("order_id", updatedOrder.OrderId), zap.String("new_status", string(updatedOrder.Status)))

	return dto.DomainOrderToProto(*updatedOrder), nil

}

// PlaceOrder implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) PlaceOrder(ctx context.Context, req *restaurantpb.PlaceOrderRequest) (*restaurantpb.PlaceOrderResponse, error) {
	if req == nil {
		return nil, domain.ErrInvalidOrderData
	}

	var orderItems []domain.OrderItem
	for _, item := range req.Items {
		orderItems = append(orderItems, domain.OrderItem{
			ItemId:   item.ItemId,
			Quantity: item.Quantity,
		})
	}

	order, err := r.restaurantUsecase.PlaceOrder(ctx, &domain.PlaceOrder{
		CustomerID:   req.CustomerId,
		RestaurantID: req.RestaurantId,
		Items:        orderItems,
	})
	if err != nil {
		return nil, err
	}

	return &restaurantpb.PlaceOrderResponse{
		OrderId:     order.OrderId,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}, nil
}

// ListRestaurants implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) ListRestaurants(
	req *restaurantpb.ListRestaurantsRequest,
	stream restaurantpb.RestaurantService_ListRestaurantsServer,
) error {
	if req == nil {
		return domain.ErrInvalidSearchData
	}

	ctx := stream.Context()

	return r.restaurantUsecase.StreamRestaurants(ctx, domain.Area{
		LatitudeMin: req.Latitude,
		LatitudeMax: req.Longitude,
		RadiusInKm:  req.RadiusKm,
	}, func(res domain.Restaurant) error {
		protoRes := &restaurantpb.Restaurant{
			RestaurantId: res.ID,
			Email:        res.Email,
			Name:         res.Name,
			Latitude:     res.Latitude,
			Longitude:    res.Longitude,
		}
		return stream.Send(protoRes)
	})
}

// Login implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) Login(ctx context.Context, req *restaurantpb.RestaurantLoginRequest) (*restaurantpb.Restaurant, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	rest, err := r.restaurantUsecase.LoginRestaurant(ctx, req.Email, req.SecretKey)
	if err != nil {
		return nil, domain.NewDomainError(domain.InvalidCredentialsMessage)
	}

	return dto.DomainRestaurantToProto(rest), nil
}

// AddMenuItem implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) AddMenuItem(ctx context.Context, req *restaurantpb.AddMenuItemRequest) (*restaurantpb.MenuItem, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	item, err := r.restaurantUsecase.AddMenuItem(ctx, req.RestaurantId, domain.MenuItem{
		ItemID:      "",
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})

	if err != nil {
		return nil, err
	}

	return &restaurantpb.MenuItem{
		ItemId:      item.ItemID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	}, nil
}

// GetRestaurant implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) GetRestaurant(ctx context.Context, req *restaurantpb.GetRestaurantRequest) (*restaurantpb.Restaurant, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	restaurant, err := r.restaurantUsecase.GetRestaurantByID(ctx, req.RestaurantId)
	if err != nil {
		return nil, err
	}

	var menuItems []*restaurantpb.MenuItem
	for _, item := range restaurant.MenuItems {
		menuItems = append(menuItems, &restaurantpb.MenuItem{
			ItemId:      item.ItemID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}

	return &restaurantpb.Restaurant{
		RestaurantId: restaurant.ID,
		Email:        restaurant.Email,
		Name:         restaurant.Name,
		Latitude:     restaurant.Latitude,
		Longitude:    restaurant.Longitude,
		Menus:        menuItems,
	}, nil
}

// RegisterRestaurant implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) RegisterRestaurant(ctx context.Context, req *restaurantpb.RegisterRestaurantRequest) (*restaurantpb.Restaurant, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	restaurant, err := r.restaurantUsecase.RegisterRestaurant(ctx, &domain.Restaurant{
		ID:        "",
		Email:     req.Email,
		Name:      req.Name,
		SecretKey: req.SecretKey,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		MenuItems: dto.ProtoRegisterMenuItemsToDomain(req.Menus),
	})
	if err != nil {
		return nil, err
	}

	return &restaurantpb.Restaurant{
		RestaurantId: restaurant.ID,
		Email:        restaurant.Email,
		Name:         restaurant.Name,
		Latitude:     restaurant.Latitude,
		Longitude:    restaurant.Longitude,
		Menus:        dto.DomainRestaurantToProto(restaurant).Menus,
	}, nil
}

// RemoveMenuItem implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) RemoveMenuItem(ctx context.Context, req *restaurantpb.RemoveMenuItemRequest) (*restaurantpb.MenuItem, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	err := r.restaurantUsecase.RemoveMenuItem(ctx, req.RestaurantId, req.ItemId)
	if err != nil {
		return nil, err
	}

	return &restaurantpb.MenuItem{
		ItemId: req.ItemId,
	}, nil
}

// UpdateMenuItem implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) UpdateMenuItem(ctx context.Context, req *restaurantpb.UpdateMenuItemRequest) (*restaurantpb.MenuItem, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	item, err := r.restaurantUsecase.UpdateMenuItem(ctx, req.RestaurantId, domain.MenuItem{
		ItemID:      req.ItemId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})

	if err != nil {
		return nil, err
	}

	return &restaurantpb.MenuItem{
		ItemId:      item.ItemID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	}, nil
}

func NewRestaurantHandler(
	server *grpc.Server, restaurantUsecase domain.RestaurantUseCase) {
	handler := &restaurantHandler{
		restaurantUsecase: restaurantUsecase,
	}
	restaurantpb.RegisterRestaurantServiceServer(server, handler)
}
