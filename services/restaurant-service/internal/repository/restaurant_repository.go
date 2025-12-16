package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
)

type restaurantRepository struct {
	db postgres.PostgresClient
}

// LoginRestaurant implements domain.RestaurantRepository.
func (r *restaurantRepository) LoginRestaurant(
	ctx context.Context,
	email string,
	secretKey string,
) (*domain.Restaurant, error) {

	query := `
		SELECT restaurant_id, email, name, latitude, longitude
		FROM restaurants
		WHERE email = $1 AND secret_key = $2
	`

	row := r.db.QueryRow(ctx, query, email, secretKey)

	var res domain.Restaurant
	err := row.Scan(
		&res.ID,
		&res.Email,
		&res.Name,
		&res.Latitude,
		&res.Longitude,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	return &res, nil
}

// CreateRestaurant implements domain.RestaurantRepository.
func (r *restaurantRepository) CreateRestaurant(
	ctx context.Context,
	restaurant *domain.Restaurant,
) (*domain.Restaurant, error) {

	query := `
		INSERT INTO restaurants (email, secret_key, name, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING restaurant_id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		restaurant.Email,
		restaurant.SecretKey,
		restaurant.Name,
		restaurant.Latitude,
		restaurant.Longitude,
	).Scan(&restaurant.ID)

	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

// GetRestaurantByID implements domain.RestaurantRepository.
func (r *restaurantRepository) GetRestaurantByID(
	ctx context.Context,
	restaurantID string,
) (*domain.Restaurant, error) {

	query := `
		SELECT restaurant_id, email, name, latitude, longitude
		FROM restaurants
		WHERE restaurant_id = $1
	`

	var res domain.Restaurant

	err := r.db.QueryRow(ctx, query, restaurantID).
		Scan(&res.ID, &res.Email, &res.Name, &res.Latitude, &res.Longitude)

	if err != nil {
		return nil, err
	}

	// Load menu items
	itemsQuery := `
		SELECT item_id, name, description, price
		FROM menu_items
		WHERE restaurant_id = $1
	`

	rows, err := r.db.Query(ctx, itemsQuery, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.MenuItem
		if err := rows.Scan(
			&item.ItemID,
			&item.Name,
			&item.Description,
			&item.Price,
		); err != nil {
			return nil, err
		}
		res.MenuItems = append(res.MenuItems, item)
	}

	return &res, nil
}

// GetRestaurants implements domain.RestaurantRepository.
func (r *restaurantRepository) StreamRestaurants(
	ctx context.Context,
	area domain.Area,
	onRow func(domain.Restaurant) error,
) error {
	query := `
		SELECT restaurant_id, email, name, latitude, longitude
		FROM restaurants
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var res domain.Restaurant
		if err := rows.Scan(&res.ID, &res.Email, &res.Name, &res.Latitude, &res.Longitude); err != nil {
			return err
		}

		if err := onRow(res); err != nil {
			return err
		}
	}

	return rows.Err()
}



// AddMenuItem implements domain.RestaurantRepository.
func (r *restaurantRepository) AddMenuItem(
	ctx context.Context,
	restaurantID string,
	item domain.MenuItem,
) (*domain.MenuItem, error) {

	query := `
		INSERT INTO menu_items (restaurant_id, name, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING item_id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		restaurantID,
		item.Name,
		item.Description,
		item.Price,
	).Scan(&item.ItemID)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// RemoveMenuItem implements domain.RestaurantRepository.
func (r *restaurantRepository) RemoveMenuItem(
	ctx context.Context,
	restaurantID string,
	itemID string,
) error {

	query := `
		DELETE FROM menu_items
		WHERE item_id = $1 AND restaurant_id = $2
	`

	affected, err := r.db.Exec(ctx, query, itemID, restaurantID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrMenuItemNotFound
	}

	return nil
}

// UpdateMenuItem implements domain.RestaurantRepository.
func (r *restaurantRepository) UpdateMenuItem(
	ctx context.Context,
	restaurantID string,
	item domain.MenuItem,
) (*domain.MenuItem, error) {

	query := `
		UPDATE menu_items
		SET name = $1, description = $2, price = $3
		WHERE item_id = $4 AND restaurant_id = $5
	`

	affected, err := r.db.Exec(
		ctx,
		query,
		item.Name,
		item.Description,
		item.Price,
		item.ItemID,
		restaurantID,
	)

	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, domain.ErrMenuItemNotFound
	}

	return &item, nil
}

// NewRestaurantRepository creates a new instance of RestaurantRepository.
func NewRestaurantRepository(db postgres.PostgresClient) domain.RestaurantRepository {
	return &restaurantRepository{db: db}
}
