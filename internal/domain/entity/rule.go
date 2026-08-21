package entity

import (
	"errors"
	"time"
)

type Rule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Version   int       `json:"version,omitempty"`
	Summary   string    `json:"summary,omitempty"`
	Content   string    `json:"content,omitempty"` // Markdown rule body
	Enabled   bool      `json:"enabled,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type RuleAssignment struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id,omitempty"`
	RuleID    string    `json:"rule_id,omitempty"`
	Enabled   bool      `json:"enabled,omitempty"`
	GrantedAt time.Time `json:"granted_at,omitempty"`
	Priority  int       `json:"priority,omitempty"` // default: 1
}

func (r Rule) Validate() error {
	if r.ID == "" {
		return errors.New("validate rule: id is required")
	}
	if r.Name == "" {
		return errors.New("validate rule: name is required")
	}
	if r.Version < 0 {
		return errors.New("validate rule: version must not be negative")
	}
	return nil
}

func (a RuleAssignment) Validate() error {
	if a.ID == "" {
		return errors.New("validate rule_assignment: id is required")
	}
	if a.AgentID == "" {
		return errors.New("validate rule_assignment: agent_id is required")
	}
	if a.RuleID == "" {
		return errors.New("validate rule_assignment: rule_id is required")
	}
	if a.Priority < 0 {
		return errors.New("validate rule_assignment: priority must not be negative")
	}
	return nil
}
