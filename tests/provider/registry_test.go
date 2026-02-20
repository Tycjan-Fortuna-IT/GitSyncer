package provider_test

import (
	"testing"

	"GitSyncer/core/provider"
)

func TestRegistryRegisterAndGetSourceControl(t *testing.T) {
	r := provider.NewProviderRegistry()

	factory := func(cfg provider.ProviderConfig) (provider.SourceControlProvider, error) {
		return &MockSourceControlProvider{Type: cfg.Type}, nil
	}

	if err := r.RegisterSourceControlProviderFactory(provider.ProviderGitHub, factory); err != nil {
		t.Fatalf("RegisterSourceControlProviderFactory failed: %v", err)
	}

	got, err := r.GetSourceControlProviderFactory(provider.ProviderGitHub)
	if err != nil {
		t.Fatalf("GetSourceControlProviderFactory failed: %v", err)
	}

	p, err := got(provider.ProviderConfig{Type: provider.ProviderGitHub})
	if err != nil {
		t.Fatalf("factory call failed: %v", err)
	}

	if p.GetProviderType() != provider.ProviderGitHub {
		t.Errorf("expected provider type %s, got %s", provider.ProviderGitHub, p.GetProviderType())
	}
}

func TestRegistryRegisterAndGetStorage(t *testing.T) {
	r := provider.NewProviderRegistry()

	factory := func(cfg provider.ProviderConfig) (provider.StorageProvider, error) {
		return &MockStorageProvider{}, nil
	}

	if err := r.RegisterStorageProviderFactory(provider.ProviderGitHub, factory); err != nil {
		t.Fatalf("RegisterStorageProviderFactory failed: %v", err)
	}

	got, err := r.GetStorageProviderFactory(provider.ProviderGitHub)
	if err != nil {
		t.Fatalf("GetStorageProviderFactory failed: %v", err)
	}

	if got == nil {
		t.Fatal("expected non-nil factory")
	}
}

func TestRegistryDuplicateSourceControlRegistration(t *testing.T) {
	r := provider.NewProviderRegistry()

	factory := func(cfg provider.ProviderConfig) (provider.SourceControlProvider, error) {
		return &MockSourceControlProvider{}, nil
	}

	if err := r.RegisterSourceControlProviderFactory(provider.ProviderGitHub, factory); err != nil {
		t.Fatalf("first RegisterSourceControlProviderFactory failed: %v", err)
	}

	err := r.RegisterSourceControlProviderFactory(provider.ProviderGitHub, factory)
	if err == nil {
		t.Fatal("expected error on duplicate source control registration, got nil")
	}
}

func TestRegistryDuplicateStorageRegistration(t *testing.T) {
	r := provider.NewProviderRegistry()

	factory := func(cfg provider.ProviderConfig) (provider.StorageProvider, error) {
		return &MockStorageProvider{}, nil
	}

	if err := r.RegisterStorageProviderFactory(provider.ProviderGitHub, factory); err != nil {
		t.Fatalf("first RegisterStorageProviderFactory failed: %v", err)
	}

	err := r.RegisterStorageProviderFactory(provider.ProviderGitHub, factory)
	if err == nil {
		t.Fatal("expected error on duplicate storage registration, got nil")
	}
}

func TestRegistryGetUnknownSourceControlProvider(t *testing.T) {
	r := provider.NewProviderRegistry()

	_, err := r.GetSourceControlProviderFactory(provider.ProviderGitHub)
	if err == nil {
		t.Fatal("expected error for unregistered source control provider, got nil")
	}
}

func TestRegistryGetUnknownStorageProvider(t *testing.T) {
	r := provider.NewProviderRegistry()

	_, err := r.GetStorageProviderFactory(provider.ProviderGitHub)
	if err == nil {
		t.Fatal("expected error for unregistered storage provider, got nil")
	}
}

func TestRegistryListSourceControl(t *testing.T) {
	r := provider.NewProviderRegistry()

	factory := func(cfg provider.ProviderConfig) (provider.SourceControlProvider, error) {
		return &MockSourceControlProvider{}, nil
	}

	_ = r.RegisterSourceControlProviderFactory(provider.ProviderGitHub, factory)
	_ = r.RegisterSourceControlProviderFactory(provider.ProviderGitLab, factory)

	types := r.ListSourceControlProviderTypes()
	if len(types) != 2 {
		t.Fatalf("expected 2 source control providers, got %d", len(types))
	}

	found := make(map[provider.ProviderType]bool)
	for _, pt := range types {
		found[pt] = true
	}

	if !found[provider.ProviderGitHub] {
		t.Error("expected GitHub in source control provider list")
	}
	if !found[provider.ProviderGitLab] {
		t.Error("expected GitLab in source control provider list")
	}
}

func TestRegistryListStorage(t *testing.T) {
	r := provider.NewProviderRegistry()

	factory := func(cfg provider.ProviderConfig) (provider.StorageProvider, error) {
		return &MockStorageProvider{}, nil
	}

	_ = r.RegisterStorageProviderFactory(provider.ProviderGitHub, factory)

	types := r.ListStorageProviderTypes()
	if len(types) != 1 {
		t.Fatalf("expected 1 storage provider, got %d", len(types))
	}
	if types[0] != provider.ProviderGitHub {
		t.Errorf("expected %s, got %s", provider.ProviderGitHub, types[0])
	}
}
