package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/orderpb"
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

func ProtoRegisterMenuItemsToDomain(items []*restaurantpb.RegisterMenuItem) []domain.MenuItem {
	var domainItems []domain.MenuItem
	for _, item := range items {
		domainItems = append(domainItems, domain.MenuItem{
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}
	return domainItems
}

func DomainOrderToProto(order domain.Order) *restaurantpb.Order {
	var orderItems []*restaurantpb.OrderItem
	for _, item := range order.Items {
		orderItems = append(orderItems, &restaurantpb.OrderItem{
			ItemId:   item.ItemId,
			Quantity: item.Quantity,
		})
	}

	var status orderpb.OrderStatus
	switch order.Status {
	case "PENDING":
		status = orderpb.OrderStatus_PENDING
	case "PREPARING":
		status = orderpb.OrderStatus_PREPARING
	case "READY":
		status = orderpb.OrderStatus_READY
	case "COMPLETED":
		status = orderpb.OrderStatus_COMPLETED
	case "CANCELLED":
		status = orderpb.OrderStatus_CANCELLED
	default:
		status = orderpb.OrderStatus_UNKNOWN
	}
	

	return &restaurantpb.Order{
		OrderId:      order.OrderId,
		CustomerId:   order.CustomerID,
		RestaurantId: order.RestaurantID,
		Items:        orderItems,
		TotalAmount:  order.TotalAmount,
		Status:       status,
	}
}	