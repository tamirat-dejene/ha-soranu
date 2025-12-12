package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
)

type userRepository struct {
	db postgres.PostgresClient
}

// GetUserPasswordHashByEmail implements domain.UserRepository.
func (u *userRepository) GetUserPasswordHashByEmail(ctx context.Context, email string) (string, error) {
	query := `
		SELECT password
		FROM users
		WHERE email = $1
	`

	var passwordHash string
	err := u.db.QueryRow(ctx, query, email).Scan(&passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", errs.OptimizedDbError(err)
	}

	return passwordHash, nil
}

// AddAddress implements domain.UserRepository.
func (u *userRepository) AddAddress(ctx context.Context, userID string, address domain.Address) (domain.Address, error) {

	query := `
		INSERT INTO addresses (
			user_id, street, city, state, postal_code, country, latitude, longitude
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING address_id, street, city, state, postal_code, country, latitude, longitude, created_at
	`

	var created domain.Address
	err := u.db.QueryRow(ctx, query,
		userID,
		address.Street,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.Latitude,
		address.Longitude,
	).Scan(
		&created.AddressID,
		&created.Street,
		&created.City,
		&created.State,
		&created.PostalCode,
		&created.Country,
		&created.Latitude,
		&created.Longitude,
		&created.CreatedAt,
	)

	if err != nil {
		return domain.Address{}, errs.OptimizedDbError(err)
	}

	return created, nil
}

// AddPhoneNumber implements domain.UserRepository.
func (u *userRepository) AddPhoneNumber(ctx context.Context, userID string, phone string) error {
	query := `
		UPDATE users
		SET phone_number = $1
		WHERE user_id = $2
	`

	rows_affected, err := u.db.Exec(ctx, query, phone, userID)
	if err != nil {
		return errs.OptimizedDbError(err)
	}

	if rows_affected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CreateUser implements domain.UserRepository.
func (u *userRepository) CreateUser(ctx context.Context, user *domain.UserRegister) (*domain.User, error) {
	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Signal for successful execution
	success := false

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if !success {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := `
		INSERT INTO users (email, username, phone_number, password)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, email, username, phone_number
	`

	var createdUser domain.User
	err = tx.QueryRow(ctx, query,
		user.Email,
		user.Username,
		user.PhoneNumber,
		user.Password,
	).Scan(
		&createdUser.UserID,
		&createdUser.Email,
		&createdUser.Username,
		&createdUser.PhoneNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", errs.OptimizedDbError(err))
	}

	createdUser.Addresses = []domain.Address{}

	success = true
	return &createdUser, nil
}

func (u *userRepository) GetAddresses(ctx context.Context, userID string) ([]domain.Address, error) {
	query := `
		SELECT address_id, street, city, state, postal_code, country, latitude, longitude, created_at
		FROM addresses
		WHERE user_id = $1
	`

	rows, err := u.db.Query(ctx, query, userID)
	if err != nil {
		return nil, errs.OptimizedDbError(err)
	}
	defer rows.Close()

	var addresses []domain.Address

	for rows.Next() {
		var a domain.Address
		err := rows.Scan(
			&a.AddressID,
			&a.Street,
			&a.City,
			&a.State,
			&a.PostalCode,
			&a.Country,
			&a.Latitude,
			&a.Longitude,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, errs.OptimizedDbError(err)
		}

		addresses = append(addresses, a)
	}

	return addresses, nil
}

// GetUserByEmail implements domain.UserRepository.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT user_id, email, username, phone_number, created_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := u.db.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PhoneNumber,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.OptimizedDbError(err)
	}

	addresses, err := u.GetAddresses(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	user.Addresses = addresses

	return &user, nil
}

// GetUserByID implements domain.UserRepository.
func (u *userRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT user_id, email, username, phone_number, created_at
		FROM users
		WHERE user_id = $1
	`

	var user domain.User
	err := u.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.PhoneNumber,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.OptimizedDbError(err)
	}

	addresses, err := u.GetAddresses(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Addresses = addresses

	return &user, nil
}

// RemoveAddress implements domain.UserRepository.
func (u *userRepository) RemoveAddress(ctx context.Context, userID string, addressID string) error {
	query := `
		DELETE FROM addresses
		WHERE address_id = $1 AND user_id = $2
	`

	rows_affected, err := u.db.Exec(ctx, query, addressID, userID)
	if err != nil {
		return errs.OptimizedDbError(err)
	}

	if rows_affected == 0 {
		return errors.New("address not found")
	}

	return nil
}

// RemovePhoneNumber implements domain.UserRepository.
func (u *userRepository) RemovePhoneNumber(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET phone_number = NULL
		WHERE user_id = $1
	`

	rows_affected, err := u.db.Exec(ctx, query, userID)
	if err != nil {
		return errs.OptimizedDbError(err)
	}

	if rows_affected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdatePhoneNumber implements domain.UserRepository.
func (u *userRepository) UpdatePhoneNumber(ctx context.Context, userID string, phone string) error {
	query := `
		UPDATE users
		SET phone_number = $1
		WHERE user_id = $2
	`
	rows_affected, err := u.db.Exec(ctx, query, phone, userID)
	if err != nil {
		return errs.OptimizedDbError(err)
	}

	if rows_affected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db postgres.PostgresClient) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}
