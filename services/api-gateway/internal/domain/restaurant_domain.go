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

type OrderItem struct {
	ItemId   string `json:"item_id"`
	Quantity int32  `json:"quantity"`
}