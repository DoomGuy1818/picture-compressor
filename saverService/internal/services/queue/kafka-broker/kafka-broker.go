package kafkaBroker

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	Reader *kafka.Reader
}

func New(dsn string, topic string, groupID string) *KafkaBroker {
	return &KafkaBroker{
		Reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: []string{dsn},
				Topic:   topic,
				GroupID: groupID,
			},
		),
	}
}

func (k *KafkaBroker) Consume(ctx context.Context) ([]byte, error) {
	const op = "services.queue.kafka.broker.consume"

	m, err := k.Reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return m.Value, nil
}
