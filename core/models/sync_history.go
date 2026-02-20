package models

import "time"

// SyncHistory represents a single sync operation audit log entry.
type SyncHistory struct {
	ID           int64      `json:"id"`
	RepositoryID int64      `json:"repository_id"`
	Status       string     `json:"status"`
	StartedAt    time.Time  `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	ErrorMessage string     `json:"error_message"`
	Details      string     `json:"details"`
}
