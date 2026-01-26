package kafkaBroker

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	Writer kafka.Writer
}

func New(dsn string, topic string) *KafkaBroker {
	return &KafkaBroker{
		Writer: kafka.Writer{
			Addr:  kafka.TCP(dsn),
			Topic: topic,
		},
	}
}

func (k *KafkaBroker) SendMessage(ctx context.Context, message string) error {
	const op = "services.queue.kafka-broker.sendMessage"
	err := k.Writer.WriteMessages(
		ctx,
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
