package provider_test

import (
	"context"

	"GitSyncer/core/models"
	"GitSyncer/core/provider"
)

var (
	_ provider.SourceControlProvider = (*MockSourceControlProvider)(nil)
	_ provider.StorageProvider       = (*MockStorageProvider)(nil)
)

// MockSourceControlProvider is a configurable mock implementation of SourceControlProvider for testing.
type MockSourceControlProvider struct {
	Type           provider.ProviderType
	Caps           []provider.SourceControlProviderCapability
	AuthenticateFn func(ctx context.Context, cred *models.Credential) error
	ListReposFn    func(ctx context.Context) ([]models.Repository, error)
	CloneRepoFn    func(ctx context.Context, repo *models.Repository, destPath string) error
	PushMirrorFn   func(ctx context.Context, repo *models.Repository, remoteURL string) error
	ValidateURLFn  func(url string) bool
}

func (m *MockSourceControlProvider) Authenticate(ctx context.Context, cred *models.Credential) error {
	if m.AuthenticateFn != nil {
		return m.AuthenticateFn(ctx, cred)
	}

	return nil
}

func (m *MockSourceControlProvider) ListRepos(ctx context.Context) ([]models.Repository, error) {
	if m.ListReposFn != nil {
		return m.ListReposFn(ctx)
	}

	return nil, nil
}

func (m *MockSourceControlProvider) CloneRepo(ctx context.Context, repo *models.Repository, destPath string) error {
	if m.CloneRepoFn != nil {
		return m.CloneRepoFn(ctx, repo, destPath)
	}

	return nil
}

func (m *MockSourceControlProvider) PushMirror(ctx context.Context, repo *models.Repository, remoteURL string) error {
	if m.PushMirrorFn != nil {
		return m.PushMirrorFn(ctx, repo, remoteURL)
	}

	return nil
}

func (m *MockSourceControlProvider) ValidateURL(url string) bool {
	if m.ValidateURLFn != nil {
		return m.ValidateURLFn(url)
	}

	return false
}

func (m *MockSourceControlProvider) GetProviderType() provider.ProviderType {
	return m.Type
}

func (m *MockSourceControlProvider) Capabilities() []provider.SourceControlProviderCapability {
	return m.Caps
}

// MockStorageProvider is a configurable mock implementation of StorageProvider for testing.
type MockStorageProvider struct {
	AuthenticateFn func(ctx context.Context, cred *models.Credential) error
	UploadFn       func(ctx context.Context, localPath string, remotePath string) error
	DownloadFn     func(ctx context.Context, remotePath string, localPath string) error
	ListFn         func(ctx context.Context, prefix string) ([]provider.StorageObject, error)
	DeleteFn       func(ctx context.Context, remotePath string) error
	GetQuotaFn     func(ctx context.Context) (*provider.QuotaInfo, error)
}

func (m *MockStorageProvider) Authenticate(ctx context.Context, cred *models.Credential) error {
	if m.AuthenticateFn != nil {
		return m.AuthenticateFn(ctx, cred)
	}

	return nil
}

func (m *MockStorageProvider) Upload(ctx context.Context, localPath string, remotePath string) error {
	if m.UploadFn != nil {
		return m.UploadFn(ctx, localPath, remotePath)
	}

	return nil
}

func (m *MockStorageProvider) Download(ctx context.Context, remotePath string, localPath string) error {
	if m.DownloadFn != nil {
		return m.DownloadFn(ctx, remotePath, localPath)
	}

	return nil
}

func (m *MockStorageProvider) List(ctx context.Context, prefix string) ([]provider.StorageObject, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, prefix)
	}

	return nil, nil
}

func (m *MockStorageProvider) Delete(ctx context.Context, remotePath string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, remotePath)
	}

	return nil
}

func (m *MockStorageProvider) GetQuota(ctx context.Context) (*provider.QuotaInfo, error) {
	if m.GetQuotaFn != nil {
		return m.GetQuotaFn(ctx)
	}

	return &provider.QuotaInfo{}, nil
}
