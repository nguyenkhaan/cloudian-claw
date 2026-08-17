# Step 2 — Core Contracts and Configuration

**Knowledge depth: 5/10**  

Read [00 — Architecture Overview](../00-architecture-overview.md), then skim [04 — Gateway and Protocol](../04-gateway-protocol.md) and [06 — Store Layer and Data Model](../06-store-data-model.md). At this point, understand why GoClaw keeps transport, provider, and persistence details outside the agent core.

Create a new module separate from GoClaw. This is the *learning implementation*, not a fork that must retain GoClaw package names.

```text
agentkit/
├── cmd/agentkit/main.go
├── internal/
│   ├── app/             # dependency composition only
│   ├── agent/           # run request, loop, prompt assembly
│   ├── model/           # canonical LLM types and Provider interface
│   ├── tools/           # tool contract, registry, policy
│   ├── session/         # history and summary persistence
│   ├── memory/          # added in Steps 6–7
│   ├── runtime/         # queue and event bus
│   └── transport/http/  # HTTP handler
└── migrations/
```

## Canonical data types

Every external format is translated into these types once. The agent loop should never decode OpenAI JSON or database rows directly.

```go
// internal/model/types.go
package model

import "context"

type Message struct {
	Role       string // system, user, assistant, tool
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
	Parameters  map[string]any // JSON Schema object
}

type Usage struct { PromptTokens, CompletionTokens, TotalTokens int }

type ChatRequest struct {
	Model    string
	Messages []Message
	Tools    []ToolDefinition
	MaxTokens int
}

type ChatResponse struct {
	Content      string
	ToolCalls    []ToolCall
	FinishReason string
	Usage        Usage
}

type Provider interface {
	Name() string
	Chat(context.Context, ChatRequest) (ChatResponse, error)
}
```

## Run and persistence contracts

```go
// internal/agent/types.go
package agent

import (
	"context"
	"agentkit/internal/model"
)

type RunRequest struct {
	TenantID, AgentID, UserID, SessionKey string
	Message string
	Stream  func(string) // optional token/event sink
}

type RunResult struct {
	Content    string
	Usage      model.Usage
	Iterations int
	ToolCalls  int
}

type SessionStore interface {
	Load(ctx context.Context, key string) (history []model.Message, summary string, err error)
	Append(ctx context.Context, key string, messages ...model.Message) error
	SetSummary(ctx context.Context, key, summary string) error
}
```

Keep the `SessionStore` interface close to the consumer. GoClaw follows the same pattern at much larger scope: `internal/store` defines focused interfaces, then `pg` and `sqlitestore` implement them. It avoids coupling the agent to `database/sql`.

## Composition root

Only `internal/app` may know concrete implementations. This makes production replacements small and test doubles easy.

```go
type App struct { Runner *agent.Runner }

func New(cfg Config, db *sql.DB, p model.Provider) (*App, error) {
	sessions := postgres.NewSessionStore(db)
	registry := tools.NewRegistry(tools.NewClock())
	registry.Register(tools.NewTime())
	return &App{Runner: agent.NewRunner(p, sessions, registry, cfg.Agent)}, nil
}
```

The next step gives the storage layer a concrete home; keep these contracts small enough that database and provider choices can change without rewriting the agent loop.
