package entity

import "time"

type Tool struct {
	ID          string    `json:"id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Version     int       `json:"version,omitempty"`
	TimeoutMs   int       `json:"timeout_ms,omitempty"` // default: 600000 (10 min)
	Enabled     bool      `json:"enabled,omitempty"`    // global toggle
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type AgentToolGrant struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id,omitempty"`
	ToolID    string    `json:"tool_id,omitempty"`
	GrantedAt time.Time `json:"granted_at,omitempty"`
}

type ToolCall struct {
	ID         string    `json:"id"`
	AgentID    string    `json:"agent_id,omitempty"`
	ToolID     string    `json:"tool_id,omitempty"`
	SessionID  string    `json:"session_id,omitempty"`
	CalledAt   time.Time `json:"called_at,omitempty"`
	ErrorStr   string    `json:"error_string,omitempty"`
	DurationMs int       `json:"duration_ms,omitempty"`
}
