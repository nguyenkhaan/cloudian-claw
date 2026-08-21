package entity

import "time"

// AgentType phân loại agent:
//   - OpenAgentType:       cấu hình mở, cho phép tinh chỉnh tham số tại phiên.
//   - PredefinedAgentType: cấu hình đóng, không cho phép tinh chỉnh tham số.
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
	AgentOwn AgentRole = "owner"
)

type Agent struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name,omitempty"`
	CreatedBy           string      `json:"created_by,omitempty"`    // owner id (single local owner)
	Provider            string      `json:"provider,omitempty"`      // denormalized default provider id
	Model               string      `json:"model,omitempty"`         // denormalized default model id
	Workspace           string      `json:"workspace,omitempty"`     // sandbox root cho file/CLI tools
	RestrictToWorkspace bool        `json:"restrict_to_workspace,omitempty"`
	ContextWindow       int         `json:"context_window,omitempty"`       // default: 200000
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

//Chia se agent cho nhieu nguoi khac su dung 
type AgentShare struct {
	UserID    string `json:"user_id,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	Role      AgentRole `json:"role,omitempty"`       // "owner", "user"
	GrantedBy string `json:"granted_by,omitempty"`
}
