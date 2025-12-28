package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tamirat-dejene/ha-soranu/services/notification-service/internal/domain"
	postgres "github.com/tamirat-dejene/ha-soranu/shared/db/pg"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

type notificationRepository struct {
	db postgres.PostgresClient
}

func NewNotificationRepository(db postgres.PostgresClient) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, recipient_id, recipient_type, order_id, title, message, is_read, type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	notificationID := uuid.NewString()
	_, err := r.db.Exec(ctx, query,
		notificationID,
		notification.RecipientID,
		notification.RecipientType,
		notification.OrderID,
		notification.Title,
		notification.Message,
		notification.IsRead,
		notification.Type,
		time.Now(),
	)

	if err != nil {
		logger.Error("failed to create notification", zap.Error(err))
		return err
	}

	logger.Info("notification created", zap.String("id", notificationID), zap.String("recipient_id", notification.RecipientID))
	return nil
}

func (r *notificationRepository) GetNotifications(ctx context.Context, recipientID string, recipientType string) ([]*domain.Notification, error) {
	logger.Info("Fetching notifications", zap.String("recipient_id", recipientID), zap.String("recipient_type", recipientType))
	query := `
		SELECT id, recipient_id, recipient_type, order_id, title, message, is_read, type, created_at
		FROM notifications
		WHERE recipient_id = $1 AND recipient_type = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, recipientID, recipientType)
	if err != nil {
		logger.Error("failed to get notifications", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(
			&n.ID,
			&n.RecipientID,
			&n.RecipientType,
			&n.OrderID,
			&n.Title,
			&n.Message,
			&n.IsRead,
			&n.Type,
			&n.CreatedAt,
		)
		if err != nil {
			logger.Error("failed to scan notification", zap.Error(err))
			return nil, err
		}
		notifications = append(notifications, &n)
	}

	return notifications, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1`

	_, err := r.db.Exec(ctx, query, notificationID)
	if err != nil {
		logger.Error("failed to mark notification as read", zap.Error(err))
		return err
	}

	return nil
}
