package sarama

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/messaging/kafka"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	p, err := sarama.NewSyncProducer(brokers, NewConfig())
	if err != nil {
		return nil, err
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Publish(ctx context.Context, msg kafka.Message) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	})

	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
