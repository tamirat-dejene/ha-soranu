package errs

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrFailedToUpdateUser = errors.New("failed to update user")
	ErrFailedToDeleteUser = errors.New("failed to delete user")
)
