package provider_test

import (
	"testing"

	"GitSyncer/core/provider"
)

func TestDetectProviderType(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected provider.ProviderType
	}{
		// GitHub
		{"github https", "https://github.com/user/repo.git", provider.ProviderGitHub},
		{"github ssh", "git@github.com:user/repo.git", provider.ProviderGitHub},
		{"github http", "http://github.com/user/repo", provider.ProviderGitHub},

		// GitLab
		{"gitlab https", "https://gitlab.com/user/repo.git", provider.ProviderGitLab},
		{"gitlab ssh", "git@gitlab.com:user/repo.git", provider.ProviderGitLab},
		{"gitlab self-hosted", "https://gitlab.example.com/user/repo.git", provider.ProviderGitLab},

		// Gitea
		{"gitea https", "https://gitea.com/user/repo.git", provider.ProviderGitea},
		{"gitea self-hosted", "https://gitea.myserver.com/user/repo.git", provider.ProviderGitea},

		// Unknown
		{"unknown provider", "https://example.com/user/repo.git", ""},
		{"empty url", "", ""},
		{"invalid url", "://invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := provider.DetectProviderType(tt.url)
			if got != tt.expected {
				t.Errorf("DetectProviderType(%q) = %q, want %q", tt.url, got, tt.expected)
			}
		})
	}
}
