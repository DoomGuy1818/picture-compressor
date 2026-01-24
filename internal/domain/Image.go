package domain

import (
	"github.com/google/uuid"
)

type Image struct {
	ID     uuid.UUID `json:"id"`
	Path   string    `json:"path"`
	Status string    `json:"status"`
}
