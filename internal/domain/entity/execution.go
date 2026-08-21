package entity

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (e Execution) Validate() error {
	if e.ID == "" {
		return errors.New("validate execution: id is required")
	}
	if e.AgentID == "" {
		return errors.New("validate execution: agent_id is required")
	}
	if e.SessionID == "" {
		return errors.New("validate execution: session_id is required")
	}
	switch e.Status {
	case ExecutionStatusCompleted, ExecutionStatusFailed, ExecutionStatusCancelled, ExecutionStatusLimitReached:
	default:
		return fmt.Errorf("validate execution: invalid status %q", e.Status)
	}
	if !e.EndedAt.IsZero() && e.EndedAt.Before(e.StartedAt) {
		return errors.New("validate execution: ended_at must not be before started_at")
	}
	return nil
}

func (ev ExecutionEvent) Validate() error {
	if ev.ID == "" {
		return errors.New("validate execution_event: id is required")
	}
	if ev.ExecutionID == "" {
		return errors.New("validate execution_event: execution_id is required")
	}
	return nil
}

func (t Trace) Validate() error {
	if t.ID == "" {
		return errors.New("validate trace: id is required")
	}
	switch t.Status {
	case TraceStatusRunning, TraceStatusCompleted, TraceStatusFailed:
	default:
		return fmt.Errorf("validate trace: invalid status %q", t.Status)
	}
	if t.TotalInputTokens < 0 {
		return errors.New("validate trace: total_input_tokens must not be negative")
	}
	if t.TotalOutputTokens < 0 {
		return errors.New("validate trace: total_output_tokens must not be negative")
	}
	if t.DurationMs < 0 {
		return errors.New("validate trace: duration_ms must not be negative")
	}
	return nil
}

func (s TraceSpan) Validate() error {
	if s.ID == "" {
		return errors.New("validate trace_span: id is required")
	}
	if s.TraceID == "" {
		return errors.New("validate trace_span: trace_id is required")
	}
	switch s.Status {
	case TraceStatusRunning, TraceStatusCompleted, TraceStatusFailed:
	default:
		return fmt.Errorf("validate trace_span: invalid status %q", s.Status)
	}
	if s.InputTokens < 0 {
		return errors.New("validate trace_span: input_tokens must not be negative")
	}
	if s.OutputTokens < 0 {
		return errors.New("validate trace_span: output_tokens must not be negative")
	}
	return nil
}
