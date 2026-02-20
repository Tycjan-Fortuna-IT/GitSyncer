package store

import (
	"database/sql"
	"fmt"
	"time"

	"GitSyncer/core/models"
)

type RepositoryStore struct {
	db *sql.DB
}

func NewRepositoryStore(db *sql.DB) *RepositoryStore {
	return &RepositoryStore{db: db}
}

func (s *RepositoryStore) Create(r *models.Repository) error {
	now := time.Now().UTC()

	result, err := s.db.Exec(
		`INSERT INTO repositories (provider_id, name, clone_url, description, is_mirror, default_branch, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		r.ProviderID, r.Name, r.CloneURL, r.Description, r.IsMirror, r.DefaultBranch, now, now,
	)
	if err != nil {
		return fmt.Errorf("RepositoryStore.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("RepositoryStore.Create: last insert id: %w", err)
	}

	r.ID = id
	r.CreatedAt = now
	r.UpdatedAt = now

	return nil
}

func (s *RepositoryStore) GetByID(id int64) (*models.Repository, error) {
	r := &models.Repository{}

	var lastSynced sql.NullTime

	err := s.db.QueryRow(
		`SELECT id, provider_id, name, clone_url, description, is_mirror, default_branch, last_synced_at, created_at, updated_at
		 FROM repositories WHERE id = ?`, id,
	).Scan(&r.ID, &r.ProviderID, &r.Name, &r.CloneURL, &r.Description, &r.IsMirror, &r.DefaultBranch, &lastSynced, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("RepositoryStore.GetByID(%d): %w", id, err)
	}

	if lastSynced.Valid {
		r.LastSyncedAt = &lastSynced.Time
	}

	return r, nil
}

func (s *RepositoryStore) List() ([]models.Repository, error) {
	rows, err := s.db.Query(
		`SELECT id, provider_id, name, clone_url, description, is_mirror, default_branch, last_synced_at, created_at, updated_at
		 FROM repositories ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("RepositoryStore.List: %w", err)
	}
	defer rows.Close()

	var repos []models.Repository

	for rows.Next() {
		var r models.Repository
		var lastSynced sql.NullTime

		if err := rows.Scan(&r.ID, &r.ProviderID, &r.Name, &r.CloneURL, &r.Description, &r.IsMirror, &r.DefaultBranch, &lastSynced, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("RepositoryStore.List: scan: %w", err)
		}

		if lastSynced.Valid {
			r.LastSyncedAt = &lastSynced.Time
		}

		repos = append(repos, r)
	}

	return repos, rows.Err()
}

func (s *RepositoryStore) Update(r *models.Repository) error {
	now := time.Now().UTC()

	result, err := s.db.Exec(
		`UPDATE repositories SET provider_id = ?, name = ?, clone_url = ?, description = ?, is_mirror = ?, default_branch = ?, last_synced_at = ?, updated_at = ?
		 WHERE id = ?`,
		r.ProviderID, r.Name, r.CloneURL, r.Description, r.IsMirror, r.DefaultBranch, r.LastSyncedAt, now, r.ID,
	)
	if err != nil {
		return fmt.Errorf("RepositoryStore.Update(%d): %w", r.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RepositoryStore.Update(%d): rows affected: %w", r.ID, err)
	}

	if rows == 0 {
		return fmt.Errorf("RepositoryStore.Update(%d): %w", r.ID, sql.ErrNoRows)
	}

	r.UpdatedAt = now

	return nil
}

func (s *RepositoryStore) Delete(id int64) error {
	result, err := s.db.Exec(`DELETE FROM repositories WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("RepositoryStore.Delete(%d): %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RepositoryStore.Delete(%d): rows affected: %w", id, err)
	}

	if rows == 0 {
		return fmt.Errorf("RepositoryStore.Delete(%d): %w", id, sql.ErrNoRows)
	}

	return nil
}
