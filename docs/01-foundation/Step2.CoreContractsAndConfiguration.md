# Step 2 — Core Contracts and Configuration

**Knowledge depth: 5/10**

This step creates provider-neutral contracts and configuration. The result is a core that can compile with fake dependencies.

## Task 1 — Define canonical model types

### Theory

Read [00 — Architecture Overview](../00-architecture-overview.md) and skim [04 — Gateway and Protocol](../04-gateway-protocol.md). External formats must be normalized once at the boundary. The agent loop must not decode OpenAI JSON, HTTP payloads, or database rows.

### Practice guide

Create `internal/model/types.go`:

```go
package model

type Message struct {
	Role       string
	Content    string
	ToolCallID string
	ToolName   string
	ToolCalls  []ToolCall
	IsError    bool
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments map[string]any
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]any
}

type Usage struct {
	PromptTokens, CompletionTokens, TotalTokens int
}

type ChatRequest struct {
	Model     string
	Messages  []Message
	Tools     []ToolDefinition
	MaxTokens int
}

type ChatResponse struct {
	Content      string
	ToolCalls    []ToolCall
	FinishReason string
	Usage        Usage
}
```

Keep the first model small. Add fields only when a later task needs durable behavior.

## Task 2 — Define provider and run contracts

### Theory

Interfaces belong close to the code that consumes them. This keeps the agent independent from a specific provider SDK or storage implementation.

### Practice guide

Add the provider interface to `internal/model`:

```go
type Provider interface {
	Name() string
	Chat(context.Context, ChatRequest) (ChatResponse, error)
}
```

Create `internal/agent/types.go`:

```go
type RunRequest struct {
	ProjectID  string
	AgentID    string
	UserID     string
	SessionKey string
	Message    string
}

type RunResult struct {
	Content    string
	Usage      model.Usage
	Iterations int
	ToolCalls  int
}
```

Do not put mutable history, counters, or cancellation on a shared runner. They belong to one `RunRequest` execution.

## Task 3 — Define the session boundary

### Theory

Read [06 — Store Layer and Data Model](../06-store-data-model.md). The agent should depend on conversation behavior, not SQL tables.

### Practice guide

Place this interface near its consumer:

```go
type SessionStore interface {
	Load(ctx context.Context, agentID, userID, key string) ([]model.Message, string, error)
	Append(ctx context.Context, agentID, userID, key string, messages ...model.Message) error
	SetSummary(ctx context.Context, agentID, userID, key, summary string) error
}
```

Scope is explicit in this teaching project. Step 14 later moves trusted identity into `context.Context` at the gateway boundary.

## Task 4 — Load and validate configuration

### Practice guide

Create `internal/app/config.go` with these groups:

```go
type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	Provider ProviderConfig
	Agent    AgentConfig
	Runtime  RuntimeConfig
}
```

At minimum, load:

- HTTP address and one demo bearer token.
- PostgreSQL URL.
- Provider base URL, API key, and model.
- Workspace root and enabled local skills.
- Context, output, iteration, tool-call, and concurrency limits.

Validate required values at startup. Do not log the provider key or bearer token. Commit an `.env.example`, not a real `.env` file.

## Task 5 — Prove dependency inversion

### Practice guide

Write two compile-time fakes:

```go
var _ model.Provider = (*fakeProvider)(nil)
var _ agent.SessionStore = (*memorySessions)(nil)
```

Add tests that construct an agent runner with both fakes. The test does not need a real model or database yet.

This step is complete when `go test ./...` passes and `internal/agent` imports neither an OpenAI SDK nor a PostgreSQL driver.
