package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/restaurantpb"
)

func DomainRestaurantToProto(r *domain.Restaurant) *restaurantpb.Restaurant {
	return &restaurantpb.Restaurant{
		RestaurantId: r.ID,
		Name:         r.Name,
		Email:        r.Email,
		Latitude: r.Latitude,
		Longitude: r.Longitude,
		Menus: 	toProtoMenuItems(r.MenuItems),
	}
}
func toProtoMenuItems(items []domain.MenuItem) []*restaurantpb.MenuItem {
	var protoItems []*restaurantpb.MenuItem
	for _, item := range items {
		protoItems = append(protoItems, &restaurantpb.MenuItem{
			ItemId:      item.ItemID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}
	return protoItems
}
