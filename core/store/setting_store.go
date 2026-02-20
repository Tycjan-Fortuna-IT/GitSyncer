package store

import (
	"database/sql"
	"fmt"
	"time"
)

type SettingStore struct {
	db *sql.DB
}

func NewSettingStore(db *sql.DB) *SettingStore {
	return &SettingStore{db: db}
}

// Get retrieves a setting value by key. Returns sql.ErrNoRows if not found.
func (s *SettingStore) Get(key string) (string, error) {
	var value string

	err := s.db.QueryRow(
		`SELECT value FROM settings WHERE key = ?`, key,
	).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("SettingStore.Get(%q): %w", key, err)
	}

	return value, nil
}

// Set inserts or updates a setting value.
func (s *SettingStore) Set(key, value string) error {
	now := time.Now().UTC()

	_, err := s.db.Exec(
		`INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		key, value, now,
	)
	if err != nil {
		return fmt.Errorf("SettingStore.Set(%q): %w", key, err)
	}

	return nil
}

// Exists returns whether a setting key exists.
func (s *SettingStore) Exists(key string) (bool, error) {
	var count int

	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM settings WHERE key = ?`, key,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("SettingStore.Exists(%q): %w", key, err)
	}

	return count > 0, nil
}

// SetInTx inserts or updates a setting value within an existing transaction.
func (s *SettingStore) SetInTx(tx *sql.Tx, key, value string) error {
	now := time.Now().UTC()

	_, err := tx.Exec(
		`INSERT INTO settings (key, value, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		key, value, now,
	)
	if err != nil {
		return fmt.Errorf("SettingStore.SetInTx(%q): %w", key, err)
	}

	return nil
}

// Delete removes a setting by key.
func (s *SettingStore) Delete(key string) error {
	_, err := s.db.Exec(`DELETE FROM settings WHERE key = ?`, key)
	if err != nil {
		return fmt.Errorf("SettingStore.Delete(%q): %w", key, err)
	}

	return nil
}
