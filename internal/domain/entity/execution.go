package entity

import (
	"encoding/json"
	"time"
)

// ExecutionStatus is the terminal status of an Agent Loop execution.
type ExecutionStatus string

const (
	ExecutionStatusCompleted    ExecutionStatus = "COMPLETED"
	ExecutionStatusFailed       ExecutionStatus = "FAILED"
	ExecutionStatusCancelled    ExecutionStatus = "CANCELLED"
	ExecutionStatusLimitReached ExecutionStatus = "LIMIT_REACHED"
)

// Execution is a single Agent Loop run.
type Execution struct {
	ID           string          `json:"id"`
	AgentID      string          `json:"agent_id,omitempty"`
	SessionID    string          `json:"session_id,omitempty"`
	Status       ExecutionStatus `json:"status,omitempty"`
	StartedAt    time.Time       `json:"started_at,omitempty"`
	EndedAt      time.Time       `json:"ended_at,omitempty"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

// ExecutionEvent is a realtime event emitted during an execution.
type ExecutionEvent struct {
	ID          string          `json:"id"`
	ExecutionID string          `json:"execution_id,omitempty"`
	SessionID   string          `json:"session_id,omitempty"`
	AgentID     string          `json:"agent_id,omitempty"`
	Type        string          `json:"type,omitempty"`
	Payload     json.RawMessage `json:"payload,omitempty"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
}

// TraceStatus is the status of a trace or telemetry span.
type TraceStatus string

const (
	TraceStatusRunning   TraceStatus = "running"
	TraceStatusCompleted TraceStatus = "completed"
	TraceStatusFailed    TraceStatus = "failed"
)

// Trace is the aggregated telemetry record of an execution (tokens, latency, error).
type Trace struct {
	ID                string      `json:"id"`
	SessionID         string      `json:"session_id,omitempty"`
	ExecutionID       string      `json:"execution_id,omitempty"`
	Status            TraceStatus `json:"status,omitempty"`
	TotalInputTokens  int64       `json:"total_input_tokens,omitempty"`
	TotalOutputTokens int64       `json:"total_output_tokens,omitempty"`
	SpanCount         int         `json:"span_count,omitempty"`
	LLMCallCount      int         `json:"llm_call_count,omitempty"`
	ToolCallCount     int         `json:"tool_call_count,omitempty"`
	DurationMs        int         `json:"duration_ms,omitempty"`
	ErrorMessage      string      `json:"error_message,omitempty"`
	StartedAt         time.Time   `json:"started_at,omitempty"`
	CompletedAt       time.Time   `json:"completed_at,omitempty"`
	CreatedAt         time.Time   `json:"created_at,omitempty"`
}

// TraceSpan is a child span of a trace (gateway, model call, tool exec...).
type TraceSpan struct {
	ID           string      `json:"id"`
	TraceID      string      `json:"trace_id,omitempty"`
	ParentSpanID string      `json:"parent_span_id,omitempty"`
	Name         string      `json:"name,omitempty"`
	Status       TraceStatus `json:"status,omitempty"`
	InputTokens  int64       `json:"input_tokens,omitempty"`
	OutputTokens int64       `json:"output_tokens,omitempty"`
	ToolCallID   string      `json:"tool_call_id,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
	StartedAt    time.Time   `json:"started_at,omitempty"`
	CompletedAt  time.Time   `json:"completed_at,omitempty"`
}
