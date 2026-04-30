package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Connect opens a Postgres connection and verifies it with a ping.
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Sensible pool defaults
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
