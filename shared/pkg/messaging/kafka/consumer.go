package kafka

import "context"

type Handler func(ctx context.Context, message *Message) error

type Consumer interface {
	Subscribe(ctx context.Context, topics []string, handler Handler) error
	Close() error
}

