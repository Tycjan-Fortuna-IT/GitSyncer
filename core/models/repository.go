package models

import "time"

// Repository represents a registered git repository linked to a provider.
type Repository struct {
	ID            int64      `json:"id"`
	ProviderID    int64      `json:"provider_id"`
	Name          string     `json:"name"`
	CloneURL      string     `json:"clone_url"`
	Description   string     `json:"description"`
	IsMirror      bool       `json:"is_mirror"`
	DefaultBranch string     `json:"default_branch"`
	LastSyncedAt  *time.Time `json:"last_synced_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
