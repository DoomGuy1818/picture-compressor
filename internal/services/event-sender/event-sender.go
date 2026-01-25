package eventSender

import (
	"context"
	"log/slog"
	"picCompressor/internal/domain/models"
	"picCompressor/internal/lib/sl"
	"time"

	"github.com/google/uuid"
)

type EventStorage interface {
	GetNewEvent() (models.Event, error)
	SetDone(ID uuid.UUID) error
}

type Sender struct {
	EventStorage EventStorage
	log          *slog.Logger
}

func New(eventStorage EventStorage, log *slog.Logger) *Sender {
	return &Sender{
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
				log.Error("failed to get new event", sl.Err(err))
				continue
			}
			//TODO: reserve some time for do job

			//TODO: send event to kafka.
			err = s.EventStorage.SetDone(event.ID)
			if err != nil {
				log.Error("failed to set event done", sl.Err(err))
				continue
			}
		}
	}()
}

func (s *Sender) SendMessage(event models.Event) error {
	//TODO: release send message in kafka
	panic("implement me")
}
