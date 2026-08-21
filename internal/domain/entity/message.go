package entity

import (
	"errors"
	"fmt"
	"time"
)

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleTool      MessageRole = "tool"
)

type MessageStatus string

const (
	MessageStatusSending   MessageStatus = "sending"
	MessageStatusCompleted MessageStatus = "completed"
	MessageStatusFailed    MessageStatus = "failed"
)

type Message struct {
	ID         string        `json:"id"`
	SessionID  string        `json:"session_id,omitempty"`
	AgentID    string        `json:"agent_id,omitempty"`
	Role       MessageRole   `json:"role,omitempty"`
	Content    string        `json:"content,omitempty"`
	Thinking   string        `json:"thinking,omitempty"`
	ToolCallID string        `json:"tool_call_id,omitempty"`
	ToolName   string        `json:"tool_name,omitempty"`
	Status     MessageStatus `json:"status,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty"`
}

func (m Message) Validate() error {
	if m.ID == "" {
		return errors.New("validate message: id is required")
	}
	if m.SessionID == "" {
		return errors.New("validate message: session_id is required")
	}
	switch m.Role {
	case MessageRoleSystem, MessageRoleUser, MessageRoleAssistant, MessageRoleTool:
	default:
		return fmt.Errorf("validate message: invalid role %q", m.Role)
	}
	switch m.Status {
	case "", MessageStatusSending, MessageStatusCompleted, MessageStatusFailed:
	default:
		return fmt.Errorf("validate message: invalid status %q", m.Status)
	}
	return nil
}
