package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	internalutil "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/util"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
)

type userRepository struct {
	db postgres.PostgresClient
}

// CreateUser implements domain.UserRepository.
func (u *userRepository) CreateUser(ctx context.Context, req domain.CreateUserRequest) (string, error) {
	// Start a transaction
	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Check if user already exists
	existingRow := tx.QueryRow(ctx, "SELECT id FROM users WHERE email=$1", req.Email)
	var existingID string
	err = existingRow.Scan(&existingID)
	if err == nil {
		// User exists
		return "", fmt.Errorf("user with email %s already exists", req.Email)
	} else if err != nil && err.Error() != "no rows in result set" {
		// Some other error
		return "", fmt.Errorf("failed to check existing user: %w", err)
	}

	// Insert new user
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	hp, _ := internalutil.HashPassword(req.Password)
	row := tx.QueryRow(ctx, query, req.Email, req.Username, hp)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// Insert addresses if provided
	for _, addr := range req.Addressess {
		addrQuery := `INSERT INTO user_addresses (user_id, address) VALUES ($1, $2)`
		if _, err := tx.Exec(ctx, addrQuery, id, addr); err != nil {
			return "", fmt.Errorf("failed to insert address: %w", err)
		}
	}

	return id, nil
}

// GetHashedPassword implements domain.UserRepository.
func (u *userRepository) GetHashedPassword(ctx context.Context, email string) (string, error) {
	row := u.db.QueryRow(ctx, "SELECT password FROM users WHERE email=$1", email)
	var hashedPassword string
	if err := row.Scan(&hashedPassword); err != nil {
		return "", fmt.Errorf("failed to get hashed password: %w", err)
	}
	return hashedPassword, nil
}

// FindByEmail implements domain.UserRepository.
func (u *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := "SELECT id, email, username FROM users WHERE email=$1"
	row := u.db.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.User{}, fmt.Errorf("user not found with email %s", email)
		}
		return domain.User{}, fmt.Errorf("failed to find user by email: %w", err)
	}

	// Address
	rows, err := u.db.Query(ctx, "SELECT address FROM user_addresses WHERE user_id=$1", user.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to query addresses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var addr string
		if err := rows.Scan(&addr); err != nil {
			return domain.User{}, fmt.Errorf("failed to scan address: %w", err)
		}
		user.Addressess = append(user.Addressess, addr)
	}

	if err := rows.Err(); err != nil {
		return domain.User{}, fmt.Errorf("error iterating addresses: %w", err)
	}

	return user, nil
}

// AddAddress implements domain.UserRepository.
func (u *userRepository) AddAddress(ctx context.Context, id string, address string) error {
	query := `INSERT INTO user_addresses (user_id, address) VALUES ($1, $2)`
	affrow, err := u.db.Exec(ctx, query, id, address)
	if err != nil {
		return fmt.Errorf("failed to add address: %w", err)
	}
	if affrow == 0 {
		return fmt.Errorf("unable to add address for user %s", id)
	}
	return nil
}

// Delete implements domain.UserRepository.
func (u *userRepository) Delete(ctx context.Context, id string) error {
	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			tx.Rollback(ctx)
		} else {
			commitErr := tx.Commit(ctx)
			if commitErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	// 1. Delete addresses
	addrQuery := `DELETE FROM user_addresses WHERE user_id=$1`
	if _, err = tx.Exec(ctx, addrQuery, id); err != nil {
		return fmt.Errorf("failed to delete addresses for user %s: %w", id, err)
	}

	// 2. Delete user
	userQuery := `DELETE FROM users WHERE id=$1`
	affrow, err := tx.Exec(ctx, userQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %w", id, err)
	}

	rowsAffected := affrow
	if rowsAffected == 0 {
		err = errors.New("user not found")
		return err // Will trigger rollback
	}

	return nil
}

// FindByID implements domain.UserRepository.
func (u *userRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	query := "SELECT id, email, username FROM users WHERE id=$1"
	row := u.db.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.User{}, fmt.Errorf("user not found with ID %s", id)
		}
		return domain.User{}, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return user, nil
}

// GetAddresses implements domain.UserRepository.
func (u *userRepository) GetAddresses(ctx context.Context, id string) ([]string, error) {
	query := "SELECT address FROM user_addresses WHERE user_id=$1"
	rows, err := u.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query addresses for user %s: %w", id, err)
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return addresses, nil
}

// RemoveAddress implements domain.UserRepository.
func (u *userRepository) RemoveAddress(ctx context.Context, id string, address string) error {
	query := `DELETE FROM user_addresses WHERE user_id=$1 AND address=$2`
	rows_affected, err := u.db.Exec(ctx, query, id, address)
	if err != nil {
		return fmt.Errorf("failed to remove address '%s' for user %s: %w", address, id, err)
	}

	if rows_affected == 0 {
		return fmt.Errorf("address '%s' not found for user %s", address, id)
	}

	return nil
}

// Update implements domain.UserRepository.
func (u *userRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		UPDATE users
		SET email=$1, username=$2
		WHERE id=$3
	`
	rows_affected, err := u.db.Exec(ctx, query, user.Email, user.Username, user.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to update user %s: %w", user.ID, err)
	}

	if rows_affected == 0 {
		return domain.User{}, fmt.Errorf("user not found with ID %s", user.ID)
	}

	return u.FindByEmail(ctx, user.Email)
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db postgres.PostgresClient) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}
