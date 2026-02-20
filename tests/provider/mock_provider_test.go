package provider_test

import (
	"context"
	"errors"
	"testing"

	"GitSyncer/core/models"
	"GitSyncer/core/provider"
)

func TestMockSourceControlProviderSatisfiesInterface(t *testing.T) {
	mock := &MockSourceControlProvider{
		Type: provider.ProviderGitHub,
		Caps: []provider.SourceControlProviderCapability{provider.CapabilitySSH, provider.CapabilityOAuth},
	}

	ctx := context.Background()

	if err := mock.Authenticate(ctx, &models.Credential{}); err != nil {
		t.Errorf("default Authenticate should return nil, got: %v", err)
	}

	repos, err := mock.ListRepos(ctx)
	if err != nil || repos != nil {
		t.Errorf("default ListRepos should return nil, nil; got: %v, %v", repos, err)
	}

	if err := mock.CloneRepo(ctx, &models.Repository{}, "/tmp"); err != nil {
		t.Errorf("default CloneRepo should return nil, got: %v", err)
	}

	if err := mock.PushMirror(ctx, &models.Repository{}, "https://example.com"); err != nil {
		t.Errorf("default PushMirror should return nil, got: %v", err)
	}

	if mock.ValidateURL("https://github.com/test") {
		t.Error("default ValidateURL should return false")
	}

	if mock.GetProviderType() != provider.ProviderGitHub {
		t.Errorf("expected %s, got %s", provider.ProviderGitHub, mock.GetProviderType())
	}

	caps := mock.Capabilities()
	if len(caps) != 2 {
		t.Fatalf("expected 2 capabilities, got %d", len(caps))
	}
}

func TestMockSourceControlProviderCustomBehavior(t *testing.T) {
	expectedErr := errors.New("auth failed")

	mock := &MockSourceControlProvider{
		Type: provider.ProviderGitLab,
		AuthenticateFn: func(ctx context.Context, cred *models.Credential) error {
			return expectedErr
		},
		ListReposFn: func(ctx context.Context) ([]models.Repository, error) {
			return []models.Repository{
				{ID: 1, Name: "test-repo"},
				{ID: 2, Name: "other-repo"},
			}, nil
		},
		ValidateURLFn: func(url string) bool {
			return url == "https://gitlab.com/test"
		},
	}

	ctx := context.Background()

	if err := mock.Authenticate(ctx, &models.Credential{}); !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}

	repos, err := mock.ListRepos(ctx)
	if err != nil {
		t.Fatalf("ListRepos error: %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(repos))
	}

	if !mock.ValidateURL("https://gitlab.com/test") {
		t.Error("expected ValidateURL to return true for matching URL")
	}
	if mock.ValidateURL("https://github.com/test") {
		t.Error("expected ValidateURL to return false for non-matching URL")
	}
}

func TestMockStorageProviderSatisfiesInterface(t *testing.T) {
	mock := &MockStorageProvider{}
	ctx := context.Background()

	if err := mock.Authenticate(ctx, &models.Credential{}); err != nil {
		t.Errorf("default Authenticate should return nil, got: %v", err)
	}

	if err := mock.Upload(ctx, "/local", "/remote"); err != nil {
		t.Errorf("default Upload should return nil, got: %v", err)
	}

	if err := mock.Download(ctx, "/remote", "/local"); err != nil {
		t.Errorf("default Download should return nil, got: %v", err)
	}

	items, err := mock.List(ctx, "prefix")
	if err != nil || items != nil {
		t.Errorf("default List should return nil, nil; got: %v, %v", items, err)
	}

	if err := mock.Delete(ctx, "/remote"); err != nil {
		t.Errorf("default Delete should return nil, got: %v", err)
	}

	quota, err := mock.GetQuota(ctx)
	if err != nil {
		t.Errorf("default GetQuota should not error, got: %v", err)
	}
	if quota == nil {
		t.Error("default GetQuota should return non-nil QuotaInfo")
	}
}

func TestMockStorageProviderCustomBehavior(t *testing.T) {
	mock := &MockStorageProvider{
		ListFn: func(ctx context.Context, prefix string) ([]provider.StorageObject, error) {
			return []provider.StorageObject{
				{Path: "backup/repo1.tar.gz", Size: 1024},
				{Path: "backup/repo2.tar.gz", Size: 2048},
			}, nil
		},
		GetQuotaFn: func(ctx context.Context) (*provider.QuotaInfo, error) {
			return &provider.QuotaInfo{
				TotalBytes: 10737418240,
				UsedBytes:  3072,
				FreeBytes:  10737415168,
			}, nil
		},
	}

	ctx := context.Background()

	items, err := mock.List(ctx, "backup/")
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Path != "backup/repo1.tar.gz" {
		t.Errorf("expected path backup/repo1.tar.gz, got %s", items[0].Path)
	}

	quota, err := mock.GetQuota(ctx)
	if err != nil {
		t.Fatalf("GetQuota error: %v", err)
	}
	if quota.TotalBytes != 10737418240 {
		t.Errorf("expected TotalBytes 10737418240, got %d", quota.TotalBytes)
	}
}
