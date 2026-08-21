package entity

import (
	"errors"
	"fmt"
	"time"
)

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

func (p Provider) Validate() error {
	if p.ID == "" {
		return errors.New("validate provider: id is required")
	}
	switch p.ProviderType {
	case "", ProviderTypeOpenAICompat:
	default:
		return fmt.Errorf("validate provider: invalid provider_type %q", p.ProviderType)
	}
	return nil
}

func (m ProviderModel) Validate() error {
	if m.ID == "" {
		return errors.New("validate provider_model: id is required")
	}
	if m.ProviderID == "" {
		return errors.New("validate provider_model: provider_id is required")
	}
	if m.ModelID == "" {
		return errors.New("validate provider_model: model_id is required")
	}
	return nil
}

func (k ProviderAPIKey) Validate() error {
	if k.ID == "" {
		return errors.New("validate provider_api_key: id is required")
	}
	if k.ProviderID == "" {
		return errors.New("validate provider_api_key: provider_id is required")
	}
	if k.KeyHash == "" {
		return errors.New("validate provider_api_key: key_hash is required")
	}
	return nil
}
