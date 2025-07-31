package postgres

import (
	"database/sql"
	"fmt"
	"url_shortener/internal/config"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	const fn = "storage.postgres.New"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const fn = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" { // PostgreSQL unique violation error code
			return fmt.Errorf("%s: duplicate entry - %w", fn, err)
		}
	}

	return nil
}
