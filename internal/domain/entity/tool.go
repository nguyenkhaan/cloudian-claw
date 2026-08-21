package entity

import (
	"errors"
	"time"
)

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

func (t Tool) Validate() error {
	if t.ID == "" {
		return errors.New("validate tool: id is required")
	}
	if t.Name == "" {
		return errors.New("validate tool: name is required")
	}
	if t.Version < 0 {
		return errors.New("validate tool: version must not be negative")
	}
	if t.TimeoutMs <= 0 {
		return errors.New("validate tool: timeout_ms must be greater than 0")
	}
	return nil
}

func (g AgentToolGrant) Validate() error {
	if g.ID == "" {
		return errors.New("validate agent_tool_grant: id is required")
	}
	if g.AgentID == "" {
		return errors.New("validate agent_tool_grant: agent_id is required")
	}
	if g.ToolID == "" {
		return errors.New("validate agent_tool_grant: tool_id is required")
	}
	return nil
}

func (c ToolCall) Validate() error {
	if c.ID == "" {
		return errors.New("validate tool_call: id is required")
	}
	if c.AgentID == "" {
		return errors.New("validate tool_call: agent_id is required")
	}
	if c.ToolID == "" {
		return errors.New("validate tool_call: tool_id is required")
	}
	if c.SessionID == "" {
		return errors.New("validate tool_call: session_id is required")
	}
	return nil
}
