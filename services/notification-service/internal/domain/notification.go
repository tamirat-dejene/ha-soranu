package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID            string
	RecipientID   string
	RecipientType string // "USER" or "RESTAURANT"
	OrderID       string
	Title         string
	Message       string
	IsRead        bool
	Type          string
	CreatedAt     time.Time
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context, recipientID string, recipientType string) ([]*Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
}

type NotificationUseCase interface {
	StartConsumer(ctx context.Context) error
	GetNotifications(ctx context.Context, recipientID string, recipientType string) ([]*Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
}
