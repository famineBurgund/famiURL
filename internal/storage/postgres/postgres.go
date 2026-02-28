package postgres

import (
	"database/sql"
	"fmt"

	"github.com/famineBurgund/famiURL/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(StoragePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", StoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS url(
        id SERIAL PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	var id int64
	err := s.db.QueryRow(`
    INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id
	`, urlToSave, alias).Scan(&id)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == "23505" {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, err
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	var res string
	err := s.db.QueryRow(`
	SELECT url FROM url WHERE alias = $1
	`, alias).Scan(&res)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == "23505" {
				return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
			}
		}
	}
	return res, err
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"
	err := s.db.QueryRow(`
	DELETE FROM url WHERE alias = $1
	`, alias).Err()
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
			}
		}
	}
	return nil
}
