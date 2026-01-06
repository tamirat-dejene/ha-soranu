package repository

import (
	"context"
	"fmt"

	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/domain"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
)

type pgRepo struct {
	db postgres.PostgresClient
}

func NewPostgresRepository(db postgres.PostgresClient) PaymentRepository {
	return &pgRepo{db: db}
}

func (r *pgRepo) Create(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, amount, currency, stripe_intent_id, client_secret, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, p.ID, p.OrderID, p.Amount, p.Currency, p.StripeIntentID, p.ClientSecret, string(p.Status))
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	return nil
}

func (r *pgRepo) UpdateStatusByIntentID(ctx context.Context, intentID string, status domain.PaymentStatus) error {
	query := `UPDATE payments SET status = $1, updated_at = NOW() WHERE stripe_intent_id = $2`
	_, err := r.db.Exec(ctx, query, string(status), intentID)
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}
	return nil
}

func (r *pgRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `SELECT id, order_id, amount, currency, stripe_intent_id, client_secret, status, created_at, updated_at
              FROM payments WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var p domain.Payment
	var status string
	if err := row.Scan(&p.ID, &p.OrderID, &p.Amount, &p.Currency, &p.StripeIntentID, &p.ClientSecret, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	p.Status = domain.PaymentStatus(status)
	return &p, nil
}
