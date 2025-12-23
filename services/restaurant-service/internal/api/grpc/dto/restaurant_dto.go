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
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
		Menus:        toProtoMenuItems(r.MenuItems),
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

	status := DomainOrderStatusToProto(order.Status)

	return &restaurantpb.Order{
		OrderId:      order.OrderId,
		CustomerId:   order.CustomerID,
		RestaurantId: order.RestaurantID,
		Items:        orderItems,
		TotalAmount:  order.TotalAmount,
		Status:       status,
	}
}

func ProtoOrderStatusToDomain(status orderpb.OrderStatus) string {
	switch status {
	case orderpb.OrderStatus_PENDING:
		return "PENDING"
	case orderpb.OrderStatus_PREPARING:
		return "PREPARING"
	case orderpb.OrderStatus_READY:
		return "READY"
	case orderpb.OrderStatus_COMPLETED:
		return "COMPLETED"
	case orderpb.OrderStatus_CANCELLED:
		return "CANCELLED"
	case orderpb.OrderStatus_SHIPPED:
		return "SHIPPED"
	default:
		return "UNKNOWN"
	}
}

func DomainOrderStatusToProto(status string) orderpb.OrderStatus {
	switch status {
	case "PENDING":
		return orderpb.OrderStatus_PENDING
	case "PREPARING":
		return orderpb.OrderStatus_PREPARING
	case "READY":
		return orderpb.OrderStatus_READY
	case "COMPLETED":
		return orderpb.OrderStatus_COMPLETED
	case "CANCELLED":
		return orderpb.OrderStatus_CANCELLED
	case "SHIPPED":
		return orderpb.OrderStatus_SHIPPED
	default:
		return orderpb.OrderStatus_UNKNOWN
	}
}
