package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tamirat-dejene/ha-soranu/services/restaurant-service/internal/domain"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

type restaurantRepository struct {
	db postgres.PostgresClient
}

// GetOrderByID implements [domain.RestaurantRepository].
func (r *restaurantRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	query := `
		SELECT order_id, restaurant_id, customer_id, total_price, status
		FROM orders
		WHERE order_id = $1
	`
	var ord domain.Order
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&ord.OrderId,
		&ord.RestaurantID,
		&ord.CustomerID,
		&ord.TotalAmount,
		&ord.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	// Load order items
	itemsQuery := `
		SELECT item_id, quantity
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ItemId, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	ord.Items = items

	return &ord, nil
}

// ShipOrder implements [domain.RestaurantRepository].
func (r *restaurantRepository) ShipOrder(ctx context.Context, restaurantID string, orderID string) (string, string, error) {
	// 1. Update order status to SHIPPED
	updateQuery := `
		UPDATE orders
		SET status = $1
		WHERE order_id = $2 AND restaurant_id = $3
		RETURNING order_id
	`

	var retOrderID string
	err := r.db.QueryRow(
		ctx,
		updateQuery,
		domain.ORDER_STATUS_SHIPPED,
		orderID,
		restaurantID,
	).Scan(&retOrderID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrOrderNotFound
		}
		return "", "", err
	}

	// 2. Select a driver (simple strategy: latest created driver)
	var driverID string
	driverQuery := `
		SELECT driver_id
		FROM drivers
		ORDER BY created_at DESC
		LIMIT 1
	`

	if derr := r.db.QueryRow(ctx, driverQuery).Scan(&driverID); derr != nil {
		// If no driver found, proceed with empty driver ID
		if !errors.Is(derr, sql.ErrNoRows) {
			return "", "", derr
		}
		driverID = ""
	}

	// 3. Return confirmation message and driver ID
	return "Order shipped successfully", driverID, nil
}

// UpdateOrderStatus implements [domain.RestaurantRepository].
func (r *restaurantRepository) UpdateOrderStatus(ctx context.Context, restaurantID string, orderID string, newStatus string) (*domain.Order, error) {
	query := `
		UPDATE orders
		SET status = $1
		WHERE order_id = $2 AND restaurant_id = $3
		RETURNING order_id, customer_id, total_price, status
	`

	var updatedOrder domain.Order

	err := r.db.QueryRow(
		ctx,
		query,
		newStatus,
		orderID,
		restaurantID,
	).Scan(
		&updatedOrder.OrderId,
		&updatedOrder.CustomerID,
		&updatedOrder.TotalAmount,
		&updatedOrder.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	updatedOrder.RestaurantID = restaurantID

	// Load order items
	itemsQuery := `
		SELECT item_id, quantity
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ItemId, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	updatedOrder.Items = items

	return &updatedOrder, nil
}

// GetOrders implements [domain.RestaurantRepository].
func (r *restaurantRepository) GetOrders(
	ctx context.Context,
	restaurantID string,
) ([]domain.Order, error) {

	// 1. Load orders
	query := `
		SELECT order_id, customer_id, total_price, status
		FROM orders
		WHERE restaurant_id = $1
		ORDER BY order_id
	`

	rows, err := r.db.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)

	for rows.Next() {
		var ord domain.Order

		err := rows.Scan(
			&ord.OrderId,
			&ord.CustomerID,
			&ord.TotalAmount,
			&ord.Status,
		)
		if err != nil {
			return nil, err
		}

		ord.RestaurantID = restaurantID

		// 2. Load order items for this order
		itemsQuery := `
			SELECT item_id, quantity
			FROM order_items
			WHERE order_id = $1
		`

		itemRows, err := r.db.Query(ctx, itemsQuery, ord.OrderId)
		if err != nil {
			return nil, err
		}

		items := make([]domain.OrderItem, 0)

		for itemRows.Next() {
			var item domain.OrderItem
			if err := itemRows.Scan(&item.ItemId, &item.Quantity); err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, item)
		}
		itemRows.Close()

		ord.Items = items
		orders = append(orders, ord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrder implements domain.RestaurantRepository.
func (r *restaurantRepository) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	query := `
		SELECT order_id, restaurant_id, customer_id, total_price, status
		FROM orders
		WHERE order_id = $1
	`
	var ord domain.Order
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&ord.OrderId,
		&ord.RestaurantID,
		&ord.CustomerID,
		&ord.TotalAmount,
		&ord.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}
	return &ord, nil
}

func calculateTotalPrice(ctx context.Context, tx postgres.Tx, restaurantID string, items []domain.OrderItem) (float64, error) {
	var total float64

	for _, it := range items {
		var price float64
		query := `
			SELECT price
			FROM menu_items
			WHERE item_id = $1 AND restaurant_id = $2
		`

		err := tx.QueryRow(ctx, query, it.ItemId, restaurantID).Scan(&price)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, errors.New("menu item not found or does not belong to restaurant")
			}
			return 0, err
		}

		total += price * float64(it.Quantity)
	}

	return total, nil
}

// PlaceOrder implements [domain.RestaurantRepository].
func (r *restaurantRepository) PlaceOrder(ctx context.Context, order *domain.PlaceOrder) (*domain.Order, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// 1. Calculate total price
	totalPrice, err := calculateTotalPrice(
		ctx,
		tx,
		order.RestaurantID,
		order.Items,
	)
	if err != nil {
		return nil, err
	}

	// 2. Create order
	var orderID string
	createOrderQuery := `
		INSERT INTO orders (customer_id, restaurant_id, total_price)
		VALUES ($1, $2, $3)
		RETURNING order_id
	`

	err = tx.QueryRow(
		ctx,
		createOrderQuery,
		order.CustomerID,
		order.RestaurantID,
		totalPrice,
	).Scan(&orderID)

	if err != nil {
		return nil, err
	}

	// 3. Insert order items
	insertItemQuery := `
		INSERT INTO order_items (order_id, item_id, quantity)
		VALUES ($1, $2, $3)
	`

	for _, it := range order.Items {
		_, err = tx.Exec(
			ctx,
			insertItemQuery,
			orderID,
			it.ItemId,
			it.Quantity,
		)
		if err != nil {
			return nil, err
		}
	}

	// 4. Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	logger.Info("placed new order", zap.String("order_id", orderID), zap.String("restaurant_id", order.RestaurantID), zap.Float64("total_price", totalPrice))

	// 5. Return created order

	return &domain.Order{
		OrderId:      orderID,
		CustomerID:   order.CustomerID,
		RestaurantID: order.RestaurantID,
		Items:        order.Items,
		TotalAmount:  totalPrice,
		Status:       "PENDING",
	}, nil
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

	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	createQuery := `
		INSERT INTO restaurants (email, secret_key, name, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING restaurant_id
	`

	err = tx.QueryRow(
		ctx,
		createQuery,
		restaurant.Email,
		restaurant.SecretKey,
		restaurant.Name,
		restaurant.Latitude,
		restaurant.Longitude,
	).Scan(&restaurant.ID)
	if err != nil {
		return nil, err
	}

	// Insert menu items if any
	insertItemQuery := `
		INSERT INTO menu_items (restaurant_id, name, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING item_id
	`

	for i := range restaurant.MenuItems {
		item := &restaurant.MenuItems[i]
		if err = tx.QueryRow(
			ctx,
			insertItemQuery,
			restaurant.ID,
			item.Name,
			item.Description,
			item.Price,
		).Scan(&item.ItemID); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
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
