package domain

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCanceled  PaymentStatus = "canceled"
)

type Payment struct {
	ID             string        `db:"id"`
	OrderID        string        `db:"order_id"`
	Amount         int64         `db:"amount"` // in smallest currency unit (e.g., cents)
	Currency       string        `db:"currency"`
	StripeIntentID string        `db:"stripe_intent_id"`
	ClientSecret   string        `db:"client_secret"`
	Status         PaymentStatus `db:"status"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
}
