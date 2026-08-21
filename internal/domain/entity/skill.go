package entity

import "time"

type Skill struct {
	ID          string    `json:"id"`
	Name        string    `json:"name,omitempty"`
	Path        string    `json:"path,omitempty"`
	FileHash    string    `json:"file_hash,omitempty"`
	Version     int       `json:"version,omitempty"`
	Description string    `json:"description,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Content     string    `json:"content,omitempty"`
	Enabled     bool      `json:"enabled,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type SkillAssignment struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id,omitempty"`
	SkillID   string    `json:"skill_id,omitempty"`
	Enabled   bool      `json:"enabled,omitempty"`
	GrantedAt time.Time `json:"granted_at,omitempty"`
}
