package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"picCompressor/internal/domain/models"
	"picCompressor/internal/lib/storage"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib" //pg driver init
)

const pictureSaved = "pictured_saved"
const ReservationTime = 120 * time.Second

type Storage struct {
	db *sql.DB
}

func New(storageURL string) (*Storage, error) {
	const op = "storage.pg.New"

	db, err := sql.Open("pgx", storageURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddPicture(id uuid.UUID, path string) error {
	const op = "storage.pg.AddPicture"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare(
		"INSERT INTO pictures (id, file_path) VALUES ($1, $2)",
	)

	if err != nil {
		return fmt.Errorf("%s: insert: %w", op, err)
	}

	eventPayload := fmt.Sprintf(`{"id": "%s", "path": "%s"}`, id.String(), path)

	_, err = stmt.Query(id, eventPayload)
	if err != nil {
		return fmt.Errorf("%s: query: %w", op, err)
	}

	if err = s.saveEvent(tx, pictureSaved, eventPayload); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func (s *Storage) GetNewEvent() (models.Event, error) {
	const op = "storage.pg.GetNewEvent"

	row := s.db.QueryRow("SELECT id, event_type, payload, reserved_to FROM pictures_events WHERE status = 'new' LIMIT 1")

	var (
		evt        storage.Event
		reservedTo sql.NullTime
	)

	err := row.Scan(&evt.ID, &evt.Type, &evt.Payload, &reservedTo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Event{}, nil
		}

		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	if reservedTo.Valid {
		evt.ReservedTo = reservedTo.Time
	} else {
		evt.ReservedTo = time.Time{} // zero-value
	}

	return models.Event{
		ID:         evt.ID,
		Type:       evt.Type,
		Payload:    evt.Payload,
		ReservedTo: evt.ReservedTo,
	}, nil
}

func (s *Storage) SetDone(ID uuid.UUID) error {
	const op = "storage.pg.SetDone"

	stmt, err := s.db.Prepare("UPDATE pictures_events SET status = 'done' WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemovePicture() {
	const op = "storage.pg.RemovePicture"
}

func (s *Storage) GetCompressedPicture() {
	const op = "storage.pg.GetCompressedPicture"
}

func (s *Storage) saveEvent(tx *sql.Tx, eventType string, payload string) error {
	const op = "storage.pg.saveEvent"
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO pictures_events (id, event_type, payload) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id, eventType, payload)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ReserveTimeForJob(ID uuid.UUID) error {
	const op = "storage.pg.ReserveTimeForJob"

	stmt, err := s.db.Prepare("UPDATE pictures_events SET reserved_to = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(time.Now().Add(ReservationTime), ID)

	return nil
}
