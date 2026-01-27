package storage

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         uuid.UUID `db:"id"`
	Type       string    `db:"event_type"`
	Payload    string    `db:"payload"`
	ReservedTo time.Time `db:"reserved_to"`
}
