package provider

import "time"

// ProviderType identifies a git hosting or storage provider.
type ProviderType string

const (
	ProviderGitHub ProviderType = "github"
	ProviderGitLab ProviderType = "gitlab"
	ProviderGitea  ProviderType = "gitea"
)

// ProviderConfig holds configuration for initializing a provider.
type ProviderConfig struct {
	Type    ProviderType      `json:"type"`
	BaseURL string            `json:"base_url"`
	Options map[string]string `json:"options,omitempty"`
}

// StorageObject represents a file or object in a storage provider.
type StorageObject struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	IsDirectory  bool      `json:"is_directory"`
}

// QuotaInfo represents storage quota information from a storage provider.
type QuotaInfo struct {
	TotalBytes int64 `json:"total_bytes"`
	UsedBytes  int64 `json:"used_bytes"`
	FreeBytes  int64 `json:"free_bytes"`
}
