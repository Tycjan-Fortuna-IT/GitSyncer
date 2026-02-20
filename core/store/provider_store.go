package store

import (
	"database/sql"
	"fmt"
	"time"

	"GitSyncer/core/models"
)

type ProviderStore struct {
	db *sql.DB
}

func NewProviderStore(db *sql.DB) *ProviderStore {
	return &ProviderStore{db: db}
}

func (s *ProviderStore) Create(p *models.Provider) error {
	now := time.Now().UTC()

	result, err := s.db.Exec(
		`INSERT INTO providers (name, type, base_url, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		p.Name, p.Type, p.BaseURL, now, now,
	)
	if err != nil {
		return fmt.Errorf("ProviderStore.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("ProviderStore.Create: last insert id: %w", err)
	}

	p.ID = id
	p.CreatedAt = now
	p.UpdatedAt = now

	return nil
}

func (s *ProviderStore) GetByID(id int64) (*models.Provider, error) {
	p := &models.Provider{}

	err := s.db.QueryRow(
		`SELECT id, name, type, base_url, created_at, updated_at
		 FROM providers WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Type, &p.BaseURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ProviderStore.GetByID(%d): %w", id, err)
	}

	return p, nil
}

func (s *ProviderStore) List() ([]models.Provider, error) {
	rows, err := s.db.Query(
		`SELECT id, name, type, base_url, created_at, updated_at
		 FROM providers ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("ProviderStore.List: %w", err)
	}
	defer rows.Close()

	var providers []models.Provider

	for rows.Next() {
		var p models.Provider

		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.BaseURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("ProviderStore.List: scan: %w", err)
		}

		providers = append(providers, p)
	}

	return providers, rows.Err()
}

func (s *ProviderStore) Update(p *models.Provider) error {
	now := time.Now().UTC()

	result, err := s.db.Exec(
		`UPDATE providers SET name = ?, type = ?, base_url = ?, updated_at = ?
		 WHERE id = ?`,
		p.Name, p.Type, p.BaseURL, now, p.ID,
	)
	if err != nil {
		return fmt.Errorf("ProviderStore.Update(%d): %w", p.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ProviderStore.Update(%d): rows affected: %w", p.ID, err)
	}

	if rows == 0 {
		return fmt.Errorf("ProviderStore.Update(%d): %w", p.ID, sql.ErrNoRows)
	}

	p.UpdatedAt = now

	return nil
}

func (s *ProviderStore) Delete(id int64) error {
	result, err := s.db.Exec(`DELETE FROM providers WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("ProviderStore.Delete(%d): %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ProviderStore.Delete(%d): rows affected: %w", id, err)
	}

	if rows == 0 {
		return fmt.Errorf("ProviderStore.Delete(%d): %w", id, sql.ErrNoRows)
	}

	return nil
}
