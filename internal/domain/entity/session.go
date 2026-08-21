package entity

import (
	"errors"
	"fmt"
	"time"
)

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
	Model           string        `json:"model,omitempty"`
	Provider        string        `json:"provider,omitempty"`
	Status          SessionStatus `json:"status,omitempty"`
	CompactionCount int           `json:"compaction_count,omitempty"`
	InputTokens     int           `json:"input_tokens,omitempty"`
	OutputTokens    int           `json:"output_tokens,omitempty"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	UpdatedAt       time.Time     `json:"updated_at,omitempty"`
}

func (s Session) Validate() error {
	if s.ID == "" {
		return errors.New("validate session: id is required")
	}
	if s.AgentID == "" {
		return errors.New("validate session: agent_id is required")
	}
	switch s.Status {
	case "", SessionStatusActive, SessionStatusArchived, SessionStatusClosed:
	default:
		return fmt.Errorf("validate session: invalid status %q", s.Status)
	}
	if s.InputTokens < 0 {
		return errors.New("validate session: input_tokens must not be negative")
	}
	if s.OutputTokens < 0 {
		return errors.New("validate session: output_tokens must not be negative")
	}
	return nil
}
