package repository

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) error
	UpdateStatusByIntentID(ctx context.Context, intentID string, status domain.PaymentStatus) error
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
}
