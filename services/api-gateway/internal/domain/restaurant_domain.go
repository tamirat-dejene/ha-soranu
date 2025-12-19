package domain

type MenuItem struct {
	ItemId      string  `json:"item_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type Restaurant struct {
	RestaurantId string     `json:"restaurant_id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	Latitude     float32    `json:"latitude"`
	Longitude    float32    `json:"longitude"`
	Menus        []MenuItem `json:"menus"`
}

type Order struct {
	OrderId      string      `json:"order_id"`
	CustomerID   string      `json:"customer_id"`
	RestaurantID string      `json:"restaurant_id"`
	Items        []OrderItem `json:"items"`
	TotalAmount  float64     `json:"total_amount"`
	Status       string      `json:"status"`
}

type OrderItem struct {
	ItemId   string `json:"item_id"`
	Quantity int32  `json:"quantity"`
}
