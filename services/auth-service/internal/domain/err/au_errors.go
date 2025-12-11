package errs

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// General errors
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternalServer = errors.New("internal server error")
	ErrDatabase       = errors.New("database error")
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict")
	ErrForbidden      = errors.New("forbidden access")

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

// User-friendly error messages
const (
	MsgInvalidRequest         = "Invalid request. Please check your input."
	MsgInternalError          = "An unexpected error occurred. Please try again later."
	MsgInvalidCredentials     = "Invalid email or password."
	MsgEmailAlreadyRegistered = "This email is already registered."
	MsgUserNotFound           = "User not found."
	MsgTokenExpired           = "Your session has expired. Please login again."
	MsgInvalidToken           = "Invalid authentication token."
	MsgUnauthorized           = "You are not authorized to perform this action."
	MsgSessionNotFound        = "Session not found. Please login again."
	MsgAddressNotFound        = "Address not found."
	MsgRefreshFailed          = "Failed to refresh token. Please login again."
)

// ToGRPCError converts internal errors to user-friendly gRPC status errors
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for specific error types
	switch {
	// Authentication errors
	case errors.Is(err, ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, MsgInvalidCredentials)
	case strings.Contains(errMsg, "bcrypt"):
		return status.Error(codes.Unauthenticated, MsgInvalidCredentials)
	case strings.Contains(errMsg, "password"):
		return status.Error(codes.Unauthenticated, MsgInvalidCredentials)
	case errors.Is(err, ErrTokenExpired):
		return status.Error(codes.Unauthenticated, MsgTokenExpired)
	case errors.Is(err, ErrInvalidToken):
		return status.Error(codes.Unauthenticated, MsgInvalidToken)
	case errors.Is(err, ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, MsgSessionNotFound)
	case errors.Is(err, ErrUnauthorized):
		return status.Error(codes.PermissionDenied, MsgUnauthorized)
	case errors.Is(err, ErrRefreshFailed):
		return status.Error(codes.Unauthenticated, MsgRefreshFailed)
	case strings.Contains(errMsg, "invalid or expired refresh token"):
		return status.Error(codes.Unauthenticated, MsgRefreshFailed)

	// User errors
	case errors.Is(err, ErrUserNotFound):
		return status.Error(codes.NotFound, MsgUserNotFound)
	case strings.Contains(errMsg, "user not found"):
		return status.Error(codes.NotFound, MsgUserNotFound)
	case errors.Is(err, ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, MsgEmailAlreadyRegistered)
	case errors.Is(err, ErrEmailAlreadyUsed):
		return status.Error(codes.AlreadyExists, MsgEmailAlreadyRegistered)
	case strings.Contains(errMsg, "already exists"):
		return status.Error(codes.AlreadyExists, MsgEmailAlreadyRegistered)
	case errors.Is(err, ErrAddressNotFound):
		return status.Error(codes.NotFound, MsgAddressNotFound)

	// Validation errors
	case errors.Is(err, ErrInvalidRequest):
		return status.Error(codes.InvalidArgument, MsgInvalidRequest)
	case errors.Is(err, ErrInvalidEmailFormat):
		return status.Error(codes.InvalidArgument, "Invalid email format.")
	case errors.Is(err, ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, "Password does not meet requirements.")

	// Database errors
	case errors.Is(err, ErrDatabase):
		return status.Error(codes.Internal, MsgInternalError)
	case strings.Contains(errMsg, "database"):
		return status.Error(codes.Internal, MsgInternalError)

	// Default to internal error for unknown errors
	default:
		return status.Error(codes.Internal, MsgInternalError)
	}
}
