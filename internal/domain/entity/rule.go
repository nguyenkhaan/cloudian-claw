package entity

import "time"

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
