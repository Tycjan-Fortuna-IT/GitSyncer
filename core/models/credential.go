package models

import "time"

// Credential represents an authentication credential for a provider
// AuthType is one of: "token", "ssh_key", "oauth"
// AuthData stores the credential value
type Credential struct {
	ID         int64     `json:"id"`
	ProviderID int64     `json:"provider_id"`
	Label      string    `json:"label"`
	AuthType   string    `json:"auth_type"`
	AuthData   string    `json:"auth_data"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
