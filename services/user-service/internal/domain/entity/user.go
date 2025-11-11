package entity

import "context"

// User represents a user entity in the system.
type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
}

// UserUsecase defines the interface for user-related business logic.
type UserUsecase interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
}