package entity

import "time"

type SessionStatus string

const (
	SessionStatusActive   SessionStatus = "active"
	SessionStatusArchived SessionStatus = "archived"
	SessionStatusClosed   SessionStatus = "closed"
)

type Session struct {
	ID              string        `json:"id"`
	AgentID         string        `json:"agent_id,omitempty"`
	UserID          string        `json:"user_id,omitempty"`
	Title           string        `json:"title,omitempty"`
	Model           string        `json:"model,omitempty"`            // denormalized
	Provider        string        `json:"provider,omitempty"`         // denormalized
	Status          SessionStatus `json:"status,omitempty"`           // active | archived | closed
	CompactionCount int           `json:"compaction_count,omitempty"`  
	InputTokens     int           `json:"input_tokens,omitempty"`
	OutputTokens    int           `json:"output_tokens,omitempty"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	UpdatedAt       time.Time     `json:"updated_at,omitempty"`
}
