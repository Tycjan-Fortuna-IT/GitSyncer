package store

import (
	"database/sql"
	"fmt"
	"time"

	"GitSyncer/core/models"
)

type CredentialStore struct {
	db *sql.DB
}

func NewCredentialStore(db *sql.DB) *CredentialStore {
	return &CredentialStore{db: db}
}

func (s *CredentialStore) Create(c *models.Credential) error {
	now := time.Now().UTC()

	result, err := s.db.Exec(
		`INSERT INTO credentials (provider_id, label, auth_type, auth_data, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		c.ProviderID, c.Label, c.AuthType, c.AuthData, now, now,
	)
	if err != nil {
		return fmt.Errorf("CredentialStore.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("CredentialStore.Create: last insert id: %w", err)
	}

	c.ID = id
	c.CreatedAt = now
	c.UpdatedAt = now

	return nil
}

func (s *CredentialStore) GetByID(id int64) (*models.Credential, error) {
	c := &models.Credential{}

	err := s.db.QueryRow(
		`SELECT id, provider_id, label, auth_type, auth_data, created_at, updated_at
		 FROM credentials WHERE id = ?`, id,
	).Scan(&c.ID, &c.ProviderID, &c.Label, &c.AuthType, &c.AuthData, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("CredentialStore.GetByID(%d): %w", id, err)
	}

	return c, nil
}

func (s *CredentialStore) GetByProviderID(providerID int64) ([]models.Credential, error) {
	rows, err := s.db.Query(
		`SELECT id, provider_id, label, auth_type, auth_data, created_at, updated_at
		 FROM credentials WHERE provider_id = ? ORDER BY id`, providerID,
	)
	if err != nil {
		return nil, fmt.Errorf("CredentialStore.GetByProviderID(%d): %w", providerID, err)
	}
	defer rows.Close()

	var creds []models.Credential

	for rows.Next() {
		var c models.Credential

		if err := rows.Scan(&c.ID, &c.ProviderID, &c.Label, &c.AuthType, &c.AuthData, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("CredentialStore.GetByProviderID(%d): scan: %w", providerID, err)
		}

		creds = append(creds, c)
	}

	return creds, rows.Err()
}

func (s *CredentialStore) List() ([]models.Credential, error) {
	rows, err := s.db.Query(
		`SELECT id, provider_id, label, auth_type, auth_data, created_at, updated_at
		 FROM credentials ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("CredentialStore.List: %w", err)
	}
	defer rows.Close()

	var creds []models.Credential

	for rows.Next() {
		var c models.Credential

		if err := rows.Scan(&c.ID, &c.ProviderID, &c.Label, &c.AuthType, &c.AuthData, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("CredentialStore.List: scan: %w", err)
		}

		creds = append(creds, c)
	}

	return creds, rows.Err()
}

func (s *CredentialStore) Update(c *models.Credential) error {
	now := time.Now().UTC()

	result, err := s.db.Exec(
		`UPDATE credentials SET provider_id = ?, label = ?, auth_type = ?, auth_data = ?, updated_at = ?
		 WHERE id = ?`,
		c.ProviderID, c.Label, c.AuthType, c.AuthData, now, c.ID,
	)
	if err != nil {
		return fmt.Errorf("CredentialStore.Update(%d): %w", c.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CredentialStore.Update(%d): rows affected: %w", c.ID, err)
	}

	if rows == 0 {
		return fmt.Errorf("CredentialStore.Update(%d): %w", c.ID, sql.ErrNoRows)
	}

	c.UpdatedAt = now

	return nil
}

func (s *CredentialStore) Delete(id int64) error {
	result, err := s.db.Exec(`DELETE FROM credentials WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("CredentialStore.Delete(%d): %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CredentialStore.Delete(%d): rows affected: %w", id, err)
	}

	if rows == 0 {
		return fmt.Errorf("CredentialStore.Delete(%d): %w", id, sql.ErrNoRows)
	}

	return nil
}
