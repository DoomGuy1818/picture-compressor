package pg

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib" //pg driver init
)

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

	_, err = tx.Exec(
		"INSERT INTO pictures (id, file_path) VALUES ($1, $2)",
		id,
		path,
	)
	if err != nil {
		return fmt.Errorf("%s: insert: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func (s *Storage) RemovePicture() {
	const op = "storage.pg.RemovePicture"
}

func (s *Storage) GetCompressedPicture() {
	const op = "storage.pg.GetCompressedPicture"
}
