package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         uuid.UUID
	Type       string
	Payload    string
	ReservedTo time.Time
}
