package provider

import (
	"fmt"
	"sync"
)

// ProviderRegistry manages registration and lookup of provider factories.
type ProviderRegistry struct {
	mu                     sync.RWMutex
	sourceControlProviders map[ProviderType]SourceControlProviderFactory
	storeProviders         map[ProviderType]StorageProviderFactory
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		sourceControlProviders: make(map[ProviderType]SourceControlProviderFactory),
		storeProviders:         make(map[ProviderType]StorageProviderFactory),
	}
}

// RegisterSourceControlProviderFactory registers a SourceControlProviderFactory for the given provider type.
func (r *ProviderRegistry) RegisterSourceControlProviderFactory(pt ProviderType, factory SourceControlProviderFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sourceControlProviders[pt]; exists {
		return fmt.Errorf("source control provider already registered: %s", pt)
	}

	r.sourceControlProviders[pt] = factory

	return nil
}

// RegisterStorageProviderFactory registers a StorageProviderFactory for the given provider type.
func (r *ProviderRegistry) RegisterStorageProviderFactory(pt ProviderType, factory StorageProviderFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storeProviders[pt]; exists {
		return fmt.Errorf("storage provider already registered: %s", pt)
	}

	r.storeProviders[pt] = factory

	return nil
}

// GetSourceControlProviderFactory returns the SourceControlProviderFactory for the given provider type.
func (r *ProviderRegistry) GetSourceControlProviderFactory(pt ProviderType) (SourceControlProviderFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.sourceControlProviders[pt]
	if !ok {
		return nil, fmt.Errorf("source control provider not found: %s", pt)
	}
	return factory, nil
}

// GetStorageProviderFactory returns the StorageProviderFactory for the given provider type.
func (r *ProviderRegistry) GetStorageProviderFactory(pt ProviderType) (StorageProviderFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.storeProviders[pt]
	if !ok {
		return nil, fmt.Errorf("storage provider not found: %s", pt)
	}
	return factory, nil
}

// ListSourceControlProviderTypes returns all registered source control provider types.
func (r *ProviderRegistry) ListSourceControlProviderTypes() []ProviderType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]ProviderType, 0, len(r.sourceControlProviders))

	for pt := range r.sourceControlProviders {
		types = append(types, pt)
	}

	return types
}

// ListStorageProviderTypes returns all registered storage provider types.
func (r *ProviderRegistry) ListStorageProviderTypes() []ProviderType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]ProviderType, 0, len(r.storeProviders))

	for pt := range r.storeProviders {
		types = append(types, pt)
	}

	return types
}
