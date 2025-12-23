package sarama

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka"
)

type Consumer struct {
	group sarama.ConsumerGroup
}

type consumerHandler struct {
	handler kafka.Handler
}

// Setup implements [sarama.ConsumerGroupHandler].
func (c *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup implements [sarama.ConsumerGroupHandler].
func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements [sarama.ConsumerGroupHandler].
func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		_ = c.handler(session.Context(), &kafka.Message{
			Topic: message.Topic,
			Key:   message.Key,
			Value: message.Value,
		})

		session.MarkMessage(message, "")
	}

	return nil
}
func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	g, err := sarama.NewConsumerGroup(brokers, groupID, NewConfig())
	if err != nil {
		return nil, err
	}

	return &Consumer{group: g}, nil
}

func (c *Consumer) Subscribe(ctx context.Context, topics []string, handler kafka.Handler) error {
	for {
		if err := c.group.Consume(ctx, topics, &consumerHandler{handler}); err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	return c.group.Close()
}
