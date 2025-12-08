package domain

import "context"

/* ----- User Entity ----- */

type User struct {
	ID       string
	Username string
	Email    string
	Password string 
	Addressess  []string
}

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
	Addressess  []string
}

type UpdateUserRequest struct {
	Username *string
	Email    *string
	Password *string
	Addressess  *[]string
}

type UserResponse struct {
	ID       string
	Username string
	Email    string
	Addressess  []string
}

type DeleteUserRequest struct {
	ID string
}

/* ----- User Usecase ----- */

type UserUsecase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (string, error)
	GetUserByID(ctx context.Context, id string) (UserResponse, error)
	GetUserHashedPassword(ctx context.Context, email string) (string, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (UserResponse, error)
	DeleteUser(ctx context.Context, id string) error

	GetUserAddresses(ctx context.Context, id string) ([]string, error)
	AddUserAddress(ctx context.Context, id string, address string) error
	RemoveUserAddress(ctx context.Context, id string, address string) error
}

/* ----- User Repository ----- */

type UserRepository interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (string, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
	GetHashedPassword(ctx context.Context, email string) (string, error)

	GetAddresses(ctx context.Context, id string) ([]string, error)
	AddAddress(ctx context.Context, id string, address string) error
	RemoveAddress(ctx context.Context, id string, address string) error
}