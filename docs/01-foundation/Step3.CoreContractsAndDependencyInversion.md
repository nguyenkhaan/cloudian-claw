# Step 3 — Core Contracts and Dependency Inversion

**Knowledge depth: 5/10**

Define provider-neutral domain contracts before connecting real infrastructure.

## Step outcome

The core compiles with fake providers and stores and imports no provider SDK or SQL driver.

## Task 1 — Define canonical model types

### Theory

- **Architectural role:** These types form the canonical model between the core and its adapters. A provider adapter translates wire types into canonical types once; agents, tools, and stores exchange only canonical types.
- **Why:** Even when field names resemble OpenAI today, their meaning belongs to the application domain. This separation allows a new provider or API version without spreading changes through the agent loop.
- **GoClaw reference:** Compare directly with [`goclaw/internal/providers/types.go`](../../goclaw/internal/providers/types.go), especially `Message`, `ToolCall`, `ToolDefinition`, `Usage`, `ChatRequest`, and `ChatResponse`. Then inspect their OpenAI translation in [`goclaw/internal/providers/adapter_openai.go`](../../goclaw/internal/providers/adapter_openai.go).

Derive these components from the data flow instead of memorizing the structs:

| Component | Question it answers | Why it is separate |
|---|---|---|
| `Message` | “What role, content, and tool relationship make up one canonical conversation item?” | It preserves conversation order and allows storage or replay without knowing a specific provider. `ToolCallID` connects a tool result to an earlier proposal, `ToolCalls` belongs to the assistant message, and `IsError` tells the model that an observation failed. |
| `ToolCall` | “What action is the model proposing, with which ID and arguments?” | It is only a **proposal**, not permission to execute. The registry in Step 11 must still validate and authorize it. The ID preserves the request/result relationship. |
| `ToolDefinition` | “Which tools does the runtime advertise to the model, and what input is valid?” | `Description` helps the model choose, while `Parameters` provides the JSON Schema used by the provider. It describes a capability but contains no execution code or permission. |
| `Usage` | “How many tokens did one model call consume?” | It supports budgets, observability, and cost tracking. Keeping it separate from text allows usage to be accumulated across multiple iterations and tool round trips. |
| `ChatRequest` | “What does the core ask one provider call to do?” | It groups the model, context messages, tool surface, and output budget into a provider-neutral input. It is not the full `RunRequest`; one run may create several `ChatRequest` values. |
| `ChatResponse` | “Did the provider return text, propose tools, or stop for another reason?” | The agent loop uses `ToolCalls` and `FinishReason` to either finish or continue the act-observe loop, and adds `Usage` to the total run usage. |

The relationship is `RunRequest → one or more ChatRequest values → ChatResponse → Message/ToolCall → tool Message → next ChatRequest`. When adding a field, ask whether it belongs to **the whole run**, **one provider call**, **one message**, or **one tool call**. The answer determines which struct should own it.

Read [00 — Architecture Overview](../00-architecture-overview.md) and skim [04 — Gateway and Protocol](../04-gateway-protocol.md). External formats must be normalized once at the boundary. The agent loop must not decode OpenAI JSON, HTTP payloads, or database rows.

### Goal

Define a shared internal language so the agent loop does not depend on OpenAI or Anthropic JSON, HTTP payloads, or SQL table shapes.

### Guide to implement

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

- **Architectural role:** `Provider.Chat` is an outbound port used by the agent. `RunRequest` and `RunResult` are the use-case contract called by transports. One run may call the provider several times and execute several tools.
- **Why:** If these layers are combined, retry logic, the tool loop, total usage, and session persistence end up inside the provider adapter. That makes fake providers harder to write and encourages transports to depend on an SDK.
- **GoClaw reference:** Compare `Provider` and `ChatRequest` in [`goclaw/internal/providers/types.go`](../../goclaw/internal/providers/types.go) with `RunRequest` and `RunResult` in [`goclaw/internal/agent/loop_types.go`](../../goclaw/internal/agent/loop_types.go), then inspect the `Loop.Run` entry point in [`goclaw/internal/agent/loop_run.go`](../../goclaw/internal/agent/loop_run.go).

Interfaces belong close to the code that consumes them. This keeps the agent independent from a specific provider SDK or storage implementation.

### Goal

Separate one LLM call from one complete agent run.

### Guide to implement

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

- **Architectural role:** `SessionStore` is a port owned by its consumer, and the PostgreSQL adapter in Step 5 implements it. The `(agent, user, key)` scope is part of the contract rather than an optional filter.
- **Why:** A behavior-focused interface keeps SQL, transactions, and schema details outside the agent. Returning the summary separately from messages also prepares for context compaction in Step 15.
- **GoClaw reference:** See how GoClaw composes smaller interfaces into `SessionStore` in [`goclaw/internal/store/session_store.go`](../../goclaw/internal/store/session_store.go), implements it for PostgreSQL in [`goclaw/internal/store/pg/sessions.go`](../../goclaw/internal/store/pg/sessions.go), and passes scope through context in [`goclaw/internal/store/context.go`](../../goclaw/internal/store/context.go).

Read [06 — Store Layer and Data Model](../06-store-data-model.md). The agent should depend on conversation behavior, not SQL tables.

### Goal

Let the agent express “load, append, and summarize this conversation” without knowing how PostgreSQL stores it.

### Guide to implement

Place this interface near its consumer:

```go
type SessionStore interface {
	Load(ctx context.Context, agentID, userID, key string) ([]model.Message, string, error)
	Append(ctx context.Context, agentID, userID, key string, messages ...model.Message) error
	SetSummary(ctx context.Context, agentID, userID, key, summary string) error
}
```

Scope is explicit in this teaching project. Step 23 later standardizes trusted identity in `context.Context` across every external boundary.

## Task 4 — Prove dependency inversion

### Theory

- **Architectural role:** A compile-time assertion verifies that a fake satisfies an interface. A constructor test verifies that the runner's dependency graph can be built entirely from test doubles.
- **Why:** An interface is useful only when its consumer can actually run with another implementation. This is an early test of unit-testability and the ability to replace providers or stores.
- **GoClaw reference:** Look for the `var _ ... = (*...)(nil)` pattern in [`goclaw/internal/providers`](../../goclaw/internal/providers) and scripted or fake dependencies in [`goclaw/internal/agent/toolloop_test.go`](../../goclaw/internal/agent/toolloop_test.go). Do not copy the size of GoClaw's fakes; model only the behavior required by the current contract.

### Goal

Prove with the compiler and tests that the core depends only on contracts, not concrete implementations.

### Guide to implement

Write two compile-time fakes:

```go
var _ model.Provider = (*fakeProvider)(nil)
var _ agent.SessionStore = (*memorySessions)(nil)
```

Add tests that construct an agent runner with both fakes. The test does not need a real model or database yet.

This step is complete when `go test ./...` passes and `internal/agent` imports neither an OpenAI SDK nor a PostgreSQL driver.
