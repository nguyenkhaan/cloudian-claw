package entity

import (
	"errors"
)

// Usign for configuration later. :))
type Config struct {
	ID                     string  `json:"id"`
	ServerHost             string  `json:"server_host,omitempty"`
	ServerPort             int     `json:"server_port,omitempty"`
	PostgresURL            string  `json:"postgres_url,omitempty"`
	DefaultProvider        string  `json:"default_provider,omitempty"`
	DefaultModel           string  `json:"default_model,omitempty"`
	AITemperature          float64 `json:"ai_temperature,omitempty"`
	AIMaxTokens            int     `json:"ai_max_tokens,omitempty"`
	QuotaMaxRequestsPerDay int     `json:"quota_max_requests_per_day,omitempty"`
	GlobalToolsEnabled     bool    `json:"global_tools_enabled,omitempty"`
	MaxToolCalls           int     `json:"max_tool_calls,omitempty"`
	MaxProviderRetries     int     `json:"max_provider_retries,omitempty"`
	MaxExecutionDuration   int     `json:"max_execution_duration,omitempty"`
	MaxContextTokens       int     `json:"max_context_tokens,omitempty"`
	MaxContinuationDepth   int     `json:"max_continuation_depth,omitempty"`
}

func (c Config) Validate() error {
	if c.ID == "" {
		return errors.New("validate config: id is required")
	}
	if c.ServerPort <= 0 {
		return errors.New("validate config: server_port must be greater than 0")
	}
	if c.PostgresURL == "" {
		return errors.New("validate config: postgres_url is required")
	}
	if c.DefaultProvider == "" {
		return errors.New("validate config: default_provider is required")
	}
	if c.DefaultModel == "" {
		return errors.New("validate config: default_model is required")
	}
	if c.AIMaxTokens <= 0 {
		return errors.New("validate config: ai_max_tokens must be greater than 0")
	}
	if c.MaxToolCalls <= 0 {
		return errors.New("validate config: max_tool_calls must be greater than 0")
	}
	if c.MaxProviderRetries < 0 {
		return errors.New("validate config: max_provider_retries must not be negative")
	}
	if c.MaxExecutionDuration <= 0 {
		return errors.New("validate config: max_execution_duration must be greater than 0")
	}
	if c.MaxContextTokens <= 0 {
		return errors.New("validate config: max_context_tokens must be greater than 0")
	}
	if c.MaxContinuationDepth <= 0 {
		return errors.New("validate config: max_continuation_depth must be greater than 0")
	}
	return nil
}
