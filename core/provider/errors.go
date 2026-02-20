package provider

import (
	"fmt"
	"time"
)

// AuthError represents an authentication or authorization failure.
type AuthError struct {
	Provider ProviderType
	Message  string
	Err      error
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("auth error [%s]: %s: %v", e.Provider, e.Message, e.Err)
	}

	return fmt.Sprintf("auth error [%s]: %s", e.Provider, e.Message)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// RateLimitError represents an API rate limit exceeded error.
type RateLimitError struct {
	Provider   ProviderType
	RetryAfter time.Duration
	Message    string
	Err        error
}

func (e *RateLimitError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("rate limit [%s]: %s (retry after %s): %v", e.Provider, e.Message, e.RetryAfter, e.Err)
	}

	return fmt.Sprintf("rate limit [%s]: %s (retry after %s)", e.Provider, e.Message, e.RetryAfter)
}

func (e *RateLimitError) Unwrap() error {
	return e.Err
}

// NetworkError represents a connectivity or timeout failure.
type NetworkError struct {
	Provider ProviderType
	Message  string
	Err      error
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("network error [%s]: %s: %v", e.Provider, e.Message, e.Err)
	}

	return fmt.Sprintf("network error [%s]: %s", e.Provider, e.Message)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}
