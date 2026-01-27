package eventConsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"saverService/internal/lib/sl"
	"saverService/internal/models"
	"saverService/internal/object"
	"saverService/internal/object/s3"
	queue "saverService/internal/services/queue"
	"time"
)

type Consumer struct {
	Consumer queue.Consumer
	Vault    object.SaverInVault
	log      *slog.Logger
}

func New(log *slog.Logger, consume queue.Consumer, v *s3.Minio) *Consumer {
	return &Consumer{
		Consumer: consume,
		log:      log,
		Vault:    v,
	}
}

func (s *Consumer) RunProcessingEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "services.event-sender.RunProcessingEvents"

	log := s.log.With("op", op)

	ticker := time.NewTicker(handlePeriod)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("stopping processing events")
				return
			case <-ticker.C:
			}

			path, err := s.GetMessage(ctx)
			if err != nil {
				log.Error("can't get message", "error", sl.Err(err))
			}

			err = s.Vault.PutObject(ctx, path)
			if err != nil {
				log.Error("can't put object", "error", sl.Err(err))
			}
		}
	}()
}

func (s *Consumer) GetMessage(ctx context.Context) (string, error) {
	const op = "services.event-sender.GetMessage"

	value, err := s.Consumer.Consume(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var payload models.PicturePayload

	err = json.Unmarshal(value, &payload)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(payload.Path)

	return payload.Path, nil
}
