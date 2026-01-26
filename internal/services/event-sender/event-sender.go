package eventSender

import (
	"context"
	"log/slog"
	"picCompressor/internal/domain/models"
	"picCompressor/internal/lib/sl"
	"picCompressor/internal/services/queue"
	"time"

	"github.com/google/uuid"
)

type EventStorage interface {
	GetNewEvent() (models.Event, error)
	SetDone(ID uuid.UUID) error
	ReserveTimeForJob(ID uuid.UUID) error
}

type Sender struct {
	Producer     queue.Sender
	EventStorage EventStorage
	log          *slog.Logger
}

func New(eventStorage EventStorage, log *slog.Logger, writer queue.Sender) *Sender {
	return &Sender{
		Producer:     writer,
		EventStorage: eventStorage,
		log:          log,
	}
}

func (s *Sender) RunProcessingEvents(ctx context.Context, handlePeriod time.Duration) {
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

			event, err := s.EventStorage.GetNewEvent()
			if err != nil {
				log.Error("failed to get new event", sl.Err(err))
				continue
			}

			if event.ID == uuid.Nil {
				log.Error("no new events")
				continue
			}

			if time.Now().After(event.ReservedTo) || event.ReservedTo.IsZero() {
				if err = s.EventStorage.ReserveTimeForJob(event.ID); err != nil {
					log.Error("failed to set reservation time", sl.Err(err))
					continue
				}

				if err = s.SendMessage(ctx, event); err != nil {
					log.Error("failed to send event", sl.Err(err))
					continue
				}

			} else {
				continue
			}

			err = s.EventStorage.SetDone(event.ID)
			if err != nil {
				log.Error("failed to set event done", sl.Err(err))
				continue
			}
		}
	}()
}

func (s *Sender) SendMessage(ctx context.Context, event models.Event) error {
	const op = "services.event-sender.SendMessage"

	if err := s.Producer.SendMessage(ctx, event.Payload); err != nil {
		log := s.log.With("op", op, "message")
		log.Error("failed to send message", sl.Err(err))
	}

	log := s.log.With("op", "services.event-sender.SendMessage")
	log.Info("sent message", slog.Any("event", event))

	return nil
}
