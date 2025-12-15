package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/restaurantpb"
	"google.golang.org/grpc"
)

type restaurantHandler struct {
	restaurantpb.UnimplementedRestaurantServiceServer
	restaurantUsecase domain.RestaurantUseCase
}

// Login implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) Login(ctx context.Context, req *restaurantpb.RestaurantLoginRequest) (*restaurantpb.RestaurantLoginResponse, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	rest, err := r.restaurantUsecase.LoginRestaurant(ctx, req.Email, req.SecretKey)
	if err != nil {
		return nil, domain.NewDomainError(domain.InvalidCredentialsMessage)
	}

	return &restaurantpb.RestaurantLoginResponse{
		Restaurant: dto.DomainRestaurantToProto(rest),
	}, nil
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
func (r *restaurantHandler) GetRestaurant(ctx context.Context, req *restaurantpb.GetRestaurantRequest) (*restaurantpb.GetRestaurantResponse, error) {
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

	return &restaurantpb.GetRestaurantResponse{
		Restaurant: &restaurantpb.Restaurant{
			RestaurantId: restaurant.ID,
			Email:        restaurant.Email,
			Name:         restaurant.Name,
			Latitude:     restaurant.Latitude,
			Longitude:    restaurant.Longitude,
			Menus:        menuItems,
		},
	}, nil
}

// ListRestaurants implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) ListRestaurants(ctx context.Context, req *restaurantpb.ListRestaurantsRequest) (*restaurantpb.ListRestaurantsResponse, error) {
	if req == nil {
		return nil, domain.ErrInvalidRestaurantData
	}

	restaurants, err := r.restaurantUsecase.GetRestaurants(ctx)
	if err != nil {
		return nil, err
	}

	var protoRestaurants []*restaurantpb.Restaurant
	for _, res := range restaurants {
		var menuItems []*restaurantpb.MenuItem
		for _, item := range res.MenuItems {
			menuItems = append(menuItems, &restaurantpb.MenuItem{
				ItemId:      item.ItemID,
				Name:        item.Name,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		protoRestaurants = append(protoRestaurants, &restaurantpb.Restaurant{
			RestaurantId: res.ID,
			Email:        res.Email,
			Name:         res.Name,
			Latitude:     res.Latitude,
			Longitude:    res.Longitude,
			Menus:        menuItems,
		})
	}

	return &restaurantpb.ListRestaurantsResponse{
		Restaurants: protoRestaurants,
	}, nil
}

// RegisterRestaurant implements restaurantpb.RestaurantServiceServer.
func (r *restaurantHandler) RegisterRestaurant(ctx context.Context, req *restaurantpb.RegisterRestaurantRequest) (*restaurantpb.RegisterRestaurantResponse, error) {
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
	})
	if err != nil {
		return nil, err
	}

	return &restaurantpb.RegisterRestaurantResponse{
		Restaurant: &restaurantpb.Restaurant{
			RestaurantId: restaurant.ID,
			Email:        restaurant.Email,
			Name:         restaurant.Name,
			Latitude:     restaurant.Latitude,
			Longitude:    restaurant.Longitude,
		},
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
