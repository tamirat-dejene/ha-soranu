package domain

import (
	"time"
)

// User represents the core User entity in the system.
type User struct {
	ID          string `json:"user_id"` // UUID
	Email       string `json:"email"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"` 
	Addresses   []Address `json:"addresses"`
	CreatedAt   time.Time `json:"created_at"`
}

// Address represents a user's address.
type Address struct {
	ID         string `json:"address_id"` // UUID
	UserID     string `json:"user_id"` // FK to User
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode uint32 `json:"postal_code"`
	Country    string `json:"country"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
	CreatedAt  time.Time `json:"created_at"`
}