package entity

import "time"

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
	Thinking   string        `json:"thinking,omitempty"` // reasoning content (thinking models)
	ToolCallID string        `json:"tool_call_id,omitempty"`
	ToolName   string        `json:"tool_name,omitempty"`
	Status     MessageStatus `json:"status,omitempty"`
	CreatedAt  time.Time     `json:"created_at,omitempty"`
}
