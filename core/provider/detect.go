package provider

import (
	"net/url"
	"strings"
)

// urlPatterns maps hostname substrings to provider types.
var urlPatterns = []struct {
	pattern      string
	providerType ProviderType
}{
	{"github.com", ProviderGitHub},
	{"gitlab.com", ProviderGitLab},
	{"gitlab.", ProviderGitLab},
	{"gitea.com", ProviderGitea},
	{"gitea.", ProviderGitea},
}

// DetectProviderType attempts to determine the provider type from a repository URL.
// Returns an empty ProviderType if the provider cannot be determined.
func DetectProviderType(rawURL string) ProviderType {
	if rawURL == "" {
		return ""
	}

	// Handle SSH URLs (e.g., git@github.com:user/repo.git)
	if strings.HasPrefix(rawURL, "git@") {
		host := extractSSHHost(rawURL)
		return matchHost(host)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	host := strings.ToLower(parsed.Hostname())

	return matchHost(host)
}

// extractSSHHost extracts the hostname from an SSH-style git URL.
// Input: "git@github.com:user/repo.git" -> "github.com"
func extractSSHHost(sshURL string) string {
	after := strings.TrimPrefix(sshURL, "git@")

	if idx := strings.Index(after, ":"); idx > 0 {
		return strings.ToLower(after[:idx])
	}

	return ""
}

// matchHost matches a hostname against known provider URL patterns.
func matchHost(host string) ProviderType {
	for _, p := range urlPatterns {
		if strings.Contains(host, p.pattern) {
			return p.providerType
		}
	}

	return ""
}
