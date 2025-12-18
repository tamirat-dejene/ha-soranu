package sarama

import (
	"github.com/IBM/sarama"
)

func NewConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_1_0_0

	cfg.Producer.Return.Successes = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	return cfg
}


