package models

import "time"

// SyncSchedule represents a cron-based sync configuration for a repository.
type SyncSchedule struct {
	ID           int64      `json:"id"`
	RepositoryID int64      `json:"repository_id"`
	CronExpr     string     `json:"cron_expr"`
	Enabled      bool       `json:"enabled"`
	LastRunAt    *time.Time `json:"last_run_at"`
	NextRunAt    *time.Time `json:"next_run_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
