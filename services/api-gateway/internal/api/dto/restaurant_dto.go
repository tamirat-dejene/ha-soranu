package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/restaurantpb"
)

type RestaurantLoginDTO struct {
	Email     string `json:"email" binding:"required,email"`
	SecretKey string `json:"secret_key" binding:"required"`
}

func (dto *RestaurantLoginDTO) ToProto() *restaurantpb.RestaurantLoginRequest {
	return &restaurantpb.RestaurantLoginRequest{
		Email:     dto.Email,
		SecretKey: dto.SecretKey,
	}
}

func RestaurantResponseFromProto(restaurant *restaurantpb.Restaurant) *domain.Restaurant {
	menuItms := make([]domain.MenuItem, 0, len(restaurant.Menus))
	for _, item := range restaurant.Menus {
		menuItms = append(menuItms, domain.MenuItem{
			ItemId:      item.ItemId,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}

	return &domain.Restaurant{
		RestaurantId: restaurant.RestaurantId,
		Name:         restaurant.Name,
		Email:        restaurant.Email,
		Latitude:     restaurant.Latitude,
		Longitude:    restaurant.Longitude,
		Menus:        menuItms,
	}
}

type RegisterRestaurantDTO struct {
	Email     string            `json:"email" binding:"required,email"`
	SecretKey string            `json:"secret_key" binding:"required"`
	Name      string            `json:"name" binding:"required"`
	Latitude  float32           `json:"latitude" binding:"required"`
	Longitude float32           `json:"longitude" binding:"required"`
	Menus     []domain.MenuItem `json:"menus"`
}

func (dto *RegisterRestaurantDTO) ToProto() *restaurantpb.RegisterRestaurantRequest {
	menuItems := make([]*restaurantpb.RegisterMenuItem, 0, len(dto.Menus))
	for _, item := range dto.Menus {
		menuItems = append(menuItems, &restaurantpb.RegisterMenuItem{
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}
	return &restaurantpb.RegisterRestaurantRequest{
		Email:     dto.Email,
		SecretKey: dto.SecretKey,
		Name:      dto.Name,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
		Menus:     menuItems,
	}
}

type ListRestaurantsDTO struct {
	Latitude  *float32 `json:"latitude" binding:"required"`
	Longitude *float32 `json:"longitude" binding:"required"`
	RadiusKm  *float32 `json:"radius_km" binding:"required"`
}

func (d *ListRestaurantsDTO) ToProto() *restaurantpb.ListRestaurantsRequest {
	return &restaurantpb.ListRestaurantsRequest{
		Latitude:  *d.Latitude,
		Longitude: *d.Longitude,
		RadiusKm:  *d.RadiusKm,
	}
}

type AddMenuItemDTO struct {
	RestaurantId string  `json:"restaurant_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Price        float32 `json:"price" binding:"required"`
}

func (dto *AddMenuItemDTO) ToProto() *restaurantpb.AddMenuItemRequest {
	return &restaurantpb.AddMenuItemRequest{
		RestaurantId: dto.RestaurantId,
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
	}
}

func MenuItemResponseFromProto(item *restaurantpb.MenuItem) *domain.MenuItem {
	return &domain.MenuItem{
		ItemId:      item.ItemId,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	}
}

type UpdateMenuItemDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float32 `json:"price" binding:"required"`
}

