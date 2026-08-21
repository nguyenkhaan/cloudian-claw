package entity

import (
	"errors"
	"fmt"
	"time"
)

type AgentType string

const (
	OpenAgentType       AgentType = "open"
	PredefinedAgentType AgentType = "predefined"
)

// AgentStatus là trạng thái vòng đời của agent.
type AgentStatus string

const (
	AgentStatusActive   AgentStatus = "active"
	AgentStatusArchived AgentStatus = "archived"
)

type AgentRole string

const (
	AgentUser AgentRole = "user"
	AgentOwn  AgentRole = "owner"
)

type Agent struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name,omitempty"`
	CreatedBy           string      `json:"created_by,omitempty"` // owner id (single local owner)
	Provider            string      `json:"provider,omitempty"`   // denormalized default provider id
	Model               string      `json:"model,omitempty"`      // denormalized default model id
	Workspace           string      `json:"workspace,omitempty"`  // sandbox root cho file/CLI tools
	RestrictToWorkspace bool        `json:"restrict_to_workspace,omitempty"`
	ContextWindow       int         `json:"context_window,omitempty"`      // default: 200000
	MaxToolIterations   int         `json:"max_tool_iterations,omitempty"` // default: 20 (bounded execution)
	Type                AgentType   `json:"type,omitempty"`
	IsDefault           bool        `json:"is_default,omitempty"`
	Status              AgentStatus `json:"status,omitempty"`
	CreatedAt           time.Time   `json:"created_at,omitempty"`
	UpdatedAt           time.Time   `json:"updated_at,omitempty"`
}

// Cac tap ngu cnah trong workspace duoc gan cho agent
type AgentContextFile struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id,omitempty"`
	Filename  string    `json:"filename,omitempty"`
	Workspace string    `json:"workspace,omitempty"`
	Content   string    `json:"content,omitempty"`
	Path      string    `json:"path,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Chia se agent cho nhieu nguoi khac su dung
type AgentShare struct {
	UserID    string    `json:"user_id,omitempty"`
	AgentID   string    `json:"agent_id,omitempty"`
	Role      AgentRole `json:"role,omitempty"` // "owner", "user"
	GrantedBy string    `json:"granted_by,omitempty"`
}

func (a *Agent) Validate() error {
	if a.ID == "" {
		return errors.New("validate agent: id is required")
	}
	if a.Name == "" {
		return errors.New("validate agent: name is required")
	}
	if a.CreatedBy == "" {
		return errors.New("validate agent: created_by is required")
	}
	if a.ContextWindow <= 0 {
		return errors.New("validate agent: context window must be greater than 0")
	}
	if a.MaxToolIterations <= 0 {
		return errors.New("validate agent: max_tool_iterations must be greater than 0")
	}
	switch a.Type {
	case "", OpenAgentType, PredefinedAgentType:
	default:
		return fmt.Errorf("validate agent: invalid type %q", a.Type)
	}
	if a.Status != "" {
		switch a.Status {
		case AgentStatusActive, AgentStatusArchived:
		default:
			return fmt.Errorf("validate agent: invalid status %q", a.Status)
		}
	}
	if a.RestrictToWorkspace && a.Workspace == "" {
		return errors.New("validate agent: workspace is required when restrict_to_workspace is true")
	}
	return nil
}

func (f AgentContextFile) Validate() error {
	if f.ID == "" {
		return errors.New("validate agent_context_file: id is required")
	}
	if f.AgentID == "" {
		return errors.New("validate agent_context_file: agent_id is required")
	}
	return nil
}

func (s AgentShare) Validate() error {
	if s.UserID == "" {
		return errors.New("validate agent_share: user_id is required")
	}
	if s.AgentID == "" {
		return errors.New("validate agent_share: agent_id is required")
	}
	switch s.Role {
	case AgentUser, AgentOwn:
	default:
		return fmt.Errorf("validate agent_share: invalid role %q", s.Role)
	}
	return nil
}
