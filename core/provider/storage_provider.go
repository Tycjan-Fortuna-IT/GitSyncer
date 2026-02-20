package provider

import (
	"context"

	"GitSyncer/core/models"
)

// StorageProvider defines the interface for interacting with a cloud storage provider.
type StorageProvider interface {
	// Authenticate validates the credential and establishes an authenticated session.
	Authenticate(ctx context.Context, cred *models.Credential) error

	// Upload sends a local file to the storage provider at the given remote path.
	Upload(ctx context.Context, localPath string, remotePath string) error

	// Download retrieves a file from the storage provider to the given local path.
	Download(ctx context.Context, remotePath string, localPath string) error

	// List returns storage objects matching the given path prefix.
	List(ctx context.Context, prefix string) ([]StorageObject, error)

	// Delete removes the object at the given remote path.
	Delete(ctx context.Context, remotePath string) error

	// GetQuota returns the storage quota information.
	GetQuota(ctx context.Context) (*QuotaInfo, error)
}

// StorageProviderFactory creates a new StorageProvider from the given configuration.
type StorageProviderFactory func(cfg ProviderConfig) (StorageProvider, error)
