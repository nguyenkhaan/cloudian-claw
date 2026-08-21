package entity

import "time"

type ProviderType string

const (
	ProviderTypeOpenAICompat ProviderType = "openai_compat"
)

// LLM Provider :v
type Provider struct {
	ID           string       `json:"id"`
	DisplayName  string       `json:"display_name,omitempty"`
	BaseURL      string       `json:"base_url,omitempty"`
	ProviderType ProviderType `json:"provider_type,omitempty"` // default: openai_compat
	Enabled      bool         `json:"enabled,omitempty"`
}

type ProviderModel struct {
	ID          string `json:"id"`
	ProviderID  string `json:"provider_id,omitempty"`
	ModelID     string `json:"model_id,omitempty"` // provider-native model id
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ProviderAPIKey struct {
	ID         string    `json:"id"`
	ProviderID string    `json:"provider_id,omitempty"`
	KeyHash    string    `json:"key_hash,omitempty"` // store hash, never plaintext
	Label      string    `json:"label,omitempty"`
	IsActive   bool      `json:"is_active,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
