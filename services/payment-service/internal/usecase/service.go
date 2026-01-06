package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/repository"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/events"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka/sarama"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/orderpb"
	"go.uber.org/zap"
)

type Service interface {
	CreatePaymentIntent(ctx context.Context, orderID string, amount int64, currency string) (*domain.Payment, error)
	HandleSucceeded(ctx context.Context, intent *stripe.PaymentIntent) error
	HandleFailed(ctx context.Context, intent *stripe.PaymentIntent) error
	GetPayment(ctx context.Context, id string) (*domain.Payment, error)
}

type service struct {
	repo      repository.PaymentRepository
	timeout   time.Duration
	publisher events.EventPublisher
}

func New(repo repository.PaymentRepository, timeout time.Duration, kafkaBroker string) (Service, error) {
	// Initialize Kafka publisher (using sarama producer)
	producer, err := sarama.NewProducer([]string{kafkaBroker})
	if err != nil {
		return nil, fmt.Errorf("kafka producer: %w", err)
	}
	pub := events.NewEventPublisher(producer)
	return &service{repo: repo, timeout: timeout, publisher: pub}, nil
}

func (s *service) CreatePaymentIntent(ctx context.Context, orderID string, amount int64, currency string) (*domain.Payment, error) {
	if orderID == "" || amount <= 0 || currency == "" {
		return nil, errors.New("invalid input")
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}
	params.AddMetadata("order_id", orderID)
	params.AutomaticPaymentMethods = &stripe.PaymentIntentAutomaticPaymentMethodsParams{Enabled: stripe.Bool(true)}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("create stripe intent: %w", err)
	}

	p := &domain.Payment{
		ID:             uuid.NewString(),
		OrderID:        orderID,
		Amount:         amount,
		Currency:       currency,
		StripeIntentID: pi.ID,
		ClientSecret:   pi.ClientSecret,
		Status:         domain.PaymentStatusPending,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("persist payment: %w", err)
	}
	return p, nil
}

func (s *service) HandleSucceeded(ctx context.Context, intent *stripe.PaymentIntent) error {
	if err := s.repo.UpdateStatusByIntentID(ctx, intent.ID, domain.PaymentStatusSucceeded); err != nil {
		return err
	}
	// Publish order status updated event
	orderID := intent.Metadata["order_id"]
	if orderID == "" {
		logger.Warn("missing order_id in intent metadata", zap.String("intent_id", intent.ID))
		return nil
	}

	ev := &orderpb.OrderStatusUpdated{
		OrderId:       orderID,
		NewStatus:     orderpb.OrderStatus_COMPLETED,
		UpdatedAtUnix: time.Now().Unix(),
	}
	if err := s.publisher.PublishOrderStatusUpdated(ctx, ev); err != nil {
		logger.Error("publish status event failed", zap.Error(err))
	}
	return nil
}

func (s *service) HandleFailed(ctx context.Context, intent *stripe.PaymentIntent) error {
	if err := s.repo.UpdateStatusByIntentID(ctx, intent.ID, domain.PaymentStatusFailed); err != nil {
		return err
	}
	orderID := intent.Metadata["order_id"]
	if orderID == "" {
		return nil
	}
	ev := &orderpb.OrderStatusUpdated{
		OrderId:       orderID,
		NewStatus:     orderpb.OrderStatus_CANCELLED,
		UpdatedAtUnix: time.Now().Unix(),
	}
	if err := s.publisher.PublishOrderStatusUpdated(ctx, ev); err != nil {
		logger.Error("publish status event failed", zap.Error(err))
	}
	return nil
}

func (s *service) GetPayment(ctx context.Context, id string) (*domain.Payment, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetByID(ctx, id)
}

//
