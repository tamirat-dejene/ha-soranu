package domain

import (
	"context"
	"time"
)

type Address struct {
    AddressID  string
    Street     string
    City       string
    State      string
    PostalCode uint32
    Country    string
    Latitude   float32
    Longitude  float32
    CreatedAt  time.Time
}

type User struct {
    UserID      string
    Email       string
    Username    string
    PhoneNumber string
    Addresses   []Address
    CreatedAt   time.Time
}

type UserUseCase interface {
    GetUser(ctx context.Context, userID string) (*User, error)

    AddPhoneNumber(ctx context.Context, userID, phone string) error
    UpdatePhoneNumber(ctx context.Context, userID, phone string) error
    RemovePhoneNumber(ctx context.Context, userID string) error

    GetAddresses(ctx context.Context, userID string) ([]Address, error)
    AddAddress(ctx context.Context, userID string, address Address) (*Address, error)
    RemoveAddress(ctx context.Context, userID, addressID string) error

    // Driver
    BeDriver(ctx context.Context, userID string) (string, error)
    GetDrivers(ctx context.Context, latitude, longitude, radius float32) ([]Driver, error)
    RemoveDriver(ctx context.Context, driverID string) error
}

type Driver struct {
    DriverID string
    User  User
}
type UserRepository interface {
    CreateUser(ctx context.Context, user *UserRegister) (*User, error)
    GetUserByID(ctx context.Context, userID string) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserPasswordHashByEmail(ctx context.Context, email string) (string, error)

    UpdatePhoneNumber(ctx context.Context, userID string, phone string) error
    AddPhoneNumber(ctx context.Context, userID string, phone string) error
    RemovePhoneNumber(ctx context.Context, userID string) error

    GetAddresses(ctx context.Context, userID string) ([]Address, error)
    AddAddress(ctx context.Context, userID string, address Address) (Address, error)
    RemoveAddress(ctx context.Context, userID, addressID string) error

    // Driver
    BeDriver(ctx context.Context, userID string) (string, error)
    GetDrivers(ctx context.Context, latitude, longitude, radius float32) ([]Driver, error)
    RemoveDriver(ctx context.Context, driverID string) error
}