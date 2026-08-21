package entity

import (
	"errors"
	"time"
)

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

func (s Skill) Validate() error {
	if s.ID == "" {
		return errors.New("validate skill: id is required")
	}
	if s.Name == "" {
		return errors.New("validate skill: name is required")
	}
	if s.Version < 0 {
		return errors.New("validate skill: version must not be negative")
	}
	return nil
}

func (a SkillAssignment) Validate() error {
	if a.ID == "" {
		return errors.New("validate skill_assignment: id is required")
	}
	if a.AgentID == "" {
		return errors.New("validate skill_assignment: agent_id is required")
	}
	if a.SkillID == "" {
		return errors.New("validate skill_assignment: skill_id is required")
	}
	return nil
}
