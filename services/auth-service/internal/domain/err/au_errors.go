package errs

import "errors"

var (
	// General errors
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInternalServer     = errors.New("internal server error")
	ErrDatabase           = errors.New("database error")
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource conflict")
	ErrForbidden          = errors.New("forbidden access")
	
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrRefreshFailed      = errors.New("unable to refresh token")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrSessionNotFound    = errors.New("session not found or already invalidated")
	ErrTokenRevoked       = errors.New("token has been revoked")

	// User domain errors
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrEmailAlreadyUsed     = errors.New("email is already registered")
	ErrUsernameAlreadyUsed  = errors.New("username is already taken")
	ErrInvalidEmailFormat   = errors.New("invalid email format")
	ErrInvalidPassword      = errors.New("password does not meet requirements")
	ErrAddressNotFound      = errors.New("address not found")
	ErrAddressAlreadyExists = errors.New("address already exists")
)
