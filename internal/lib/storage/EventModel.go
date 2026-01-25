package storage

import "github.com/google/uuid"

type Event struct {
	ID      uuid.UUID `db:"id"`
	Type    string    `db:"event_type"`
	Payload string    `db:"payload"`
}
