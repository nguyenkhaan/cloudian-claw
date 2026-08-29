package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ApplicationConfig is the typed configuration required to start Cloudclaw.
type ApplicationConfig struct {
	GatewayAPIToken  string
	OpenRouterAPIKey string

	ServerHost  string
	ServerPort  int
	PostgresURL string

	DefaultProvider string
	DefaultModel    string
	AITemperature   float64
	AIMaxTokens     int

	QuotaMaxRequestsPerDay int
	GlobalToolsEnabled     bool
	MaxToolCalls           int
	MaxProviderRetries     int
	MaxExecutionDuration   time.Duration
	MaxContextTokens       int
	MaxContinuationDepth   int
}

// LoadApplicationConfig reads typed values from the process environment.
// Ham main se tien hanh load .env thay the 
func LoadApplicationConfig() (*ApplicationConfig, error) {
	serverPort, err := requiredInt("SERVER_PORT")
	if err != nil {
		return nil, err
	}
	temperature, err := requiredFloat64("AI_TEMPERATURE")
	if err != nil {
		return nil, err
	}
	maxTokens, err := requiredInt("AI_MAX_TOKENS")
	if err != nil {
		return nil, err
	}
	quota, err := requiredInt("QUOTA_MAX_REQUESTS_PER_DAY")
	if err != nil {
		return nil, err
	}
	toolsEnabled, err := requiredBool("GLOBAL_TOOLS_ENABLED")
	if err != nil {
		return nil, err
	}
	maxToolCalls, err := requiredInt("MAX_TOOL_CALLS")
	if err != nil {
		return nil, err
	}
	maxRetries, err := requiredInt("MAX_PROVIDER_RETRIES")
	if err != nil {
		return nil, err
	}
	maxDuration, err := requiredDuration("MAX_EXECUTION_DURATION")
	if err != nil {
		return nil, err
	}
	maxContextTokens, err := requiredInt("MAX_CONTEXT_TOKENS")
	if err != nil {
		return nil, err
	}
	maxDepth, err := requiredInt("MAX_CONTINUATION_DEPTH")
	if err != nil {
		return nil, err
	}

	config := &ApplicationConfig{
		GatewayAPIToken:        requiredString("GATEWAY_API_TOKEN"),
		OpenRouterAPIKey:       requiredString("OPENROUTER_API_KEY"),
		ServerHost:             requiredString("SERVER_HOST"),
		ServerPort:             serverPort,
		PostgresURL:            requiredString("POSTGRES_URL"),
		DefaultProvider:        requiredString("DEFAULT_PROVIDER"),
		DefaultModel:           requiredString("DEFAULT_MODEL"),
		AITemperature:          temperature,
		AIMaxTokens:            maxTokens,
		QuotaMaxRequestsPerDay: quota,
		GlobalToolsEnabled:     toolsEnabled,
		MaxToolCalls:           maxToolCalls,
		MaxProviderRetries:     maxRetries,
		MaxExecutionDuration:   maxDuration,
		MaxContextTokens:       maxContextTokens,
		MaxContinuationDepth:   maxDepth,
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c ApplicationConfig) Validate() error {
	if c.GatewayAPIToken == "" || c.ServerHost == "" || c.PostgresURL == "" || c.DefaultProvider == "" || c.DefaultModel == "" {
		return fmt.Errorf("validate application config: a required string value is empty")
	}
	if c.ServerPort <= 0 || c.AIMaxTokens <= 0 || c.MaxToolCalls <= 0 || c.MaxExecutionDuration <= 0 || c.MaxContextTokens <= 0 || c.MaxContinuationDepth <= 0 {
		return fmt.Errorf("validate application config: a positive value is required")
	}
	if c.MaxProviderRetries < 0 || c.QuotaMaxRequestsPerDay < 0 {
		return fmt.Errorf("validate application config: retry and quota values cannot be negative")
	}
	return nil
}

func requiredString(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func requiredInt(name string) (int, error) {
	parsed, err := strconv.Atoi(requiredString(name))
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	return parsed, nil
}

func requiredFloat64(name string) (float64, error) {
	parsed, err := strconv.ParseFloat(requiredString(name), 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	return parsed, nil
}

func requiredBool(name string) (bool, error) {
	parsed, err := strconv.ParseBool(requiredString(name))
	if err != nil {
		return false, fmt.Errorf("parse %s: %w", name, err)
	}
	return parsed, nil
}

func requiredDuration(name string) (time.Duration, error) {
	parsed, err := time.ParseDuration(requiredString(name))
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", name, err)
	}
	return parsed, nil
}
