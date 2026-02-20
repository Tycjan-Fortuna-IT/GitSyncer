package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("database.Open: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()

		return nil, fmt.Errorf("database.Open: enable WAL: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()

		return nil, fmt.Errorf("database.Open: enable foreign keys: %w", err)
	}

	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		db.Close()

		return nil, fmt.Errorf("database.Open: set busy timeout: %w", err)
	}

	return db, nil
}

// %APPDATA%/GitSyncer/gitsyncer.db this style, platform specific
func DefaultDBPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("database.DefaultDBPath: %w", err)
	}

	dir := filepath.Join(configDir, "GitSyncer")
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("database.DefaultDBPath: create dir: %w", err)
	}

	return filepath.Join(dir, "gitsyncer.db"), nil
}

func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}

	return db.Close()
}
