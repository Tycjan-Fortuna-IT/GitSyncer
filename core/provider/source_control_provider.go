package provider

import (
	"context"

	"GitSyncer/core/models"
)

// SourceControlProviderCapability represents a feature supported by a source control provider.
type SourceControlProviderCapability string

const (
	CapabilityWebhooks  SourceControlProviderCapability = "webhooks"
	CapabilitySSH       SourceControlProviderCapability = "ssh"
	CapabilityOAuth     SourceControlProviderCapability = "oauth"
	CapabilityTokenAuth SourceControlProviderCapability = "token_auth"
	CapabilityMirror    SourceControlProviderCapability = "mirror"
)

// SourceControlProvider defines the interface for interacting with a source control provider.
type SourceControlProvider interface {
	// Authenticate validates the credential and establishes an authenticated session.
	Authenticate(ctx context.Context, cred *models.Credential) error

	// ListRepos returns all repositories accessible with the current credentials.
	ListRepos(ctx context.Context) ([]models.Repository, error)

	// CloneRepo clones the repository to the specified local path.
	CloneRepo(ctx context.Context, repo *models.Repository, destPath string) error

	// PushMirror pushes a local repository as a mirror to the given remote URL.
	PushMirror(ctx context.Context, repo *models.Repository, remoteURL string) error

	// ValidateURL checks whether the given URL belongs to this provider.
	ValidateURL(url string) bool

	// GetProviderType returns the provider type identifier.
	GetProviderType() ProviderType

	// Capabilities returns the set of features supported by this provider.
	Capabilities() []SourceControlProviderCapability
}

// SourceControlProviderFactory creates a new SourceControlProvider from the given configuration.
type SourceControlProviderFactory func(cfg ProviderConfig) (SourceControlProvider, error)
