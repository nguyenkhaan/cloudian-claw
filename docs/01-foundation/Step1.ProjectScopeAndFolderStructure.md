# Step 1 — Project Scope and Folder Structure

**Knowledge depth: 4/10**

This step defines the product boundary and creates the project skeleton. Do not implement agent behavior yet.

## Task 1 — Understand the runtime path

### Theory

Read [00 — Architecture Overview](../00-architecture-overview.md). Focus on the startup path, package responsibilities, and the boundary between the agent core and GoClaw's platform features.

GoClaw is a gateway, not only an LLM wrapper. A request enters through a transport, waits in a session queue, runs through an agent pipeline, may call tools, persists durable state, and returns through the same transport.

```text
request → session queue → agent pipeline → provider/tools → store → response
```

### Practice guide

Write this runtime path in the project README. For each arrow, identify the package that will own it:

| Runtime responsibility | Project package |
|---|---|
| Process composition | `internal/app` |
| Agent execution | `internal/agent` |
| Provider-neutral types | `internal/model` |
| Tool execution and policy | `internal/tools` |
| Session persistence | `internal/session`, `internal/store/postgres` |
| Queue and events | `internal/runtime` |
| HTTP and WebSocket | `internal/transport` |

The mapping is complete when one responsibility has one clear owner.

## Task 2 — Fix the project scope

### Theory

The course builds one complete vertical slice:

```text
One Go gateway + one React web UI
One OpenAI-compatible provider
PostgreSQL with pgvector
One agent profile and per-user sessions
Conversation memory and a small tool set
HTTP API + WebSocket streaming
Structured logs and a simple run timeline
```

Channels, teams, multiple providers, skill publishing, cron/heartbeat, and multi-tenancy are useful GoClaw architecture topics. They are not part of the first implementation.

### Practice guide

Add a short `Scope` section to the project README with two lists:

- `Build now`: the seven items above.
- `Later extensions`: channels, teams, provider routing, skill publishing, autonomous schedules, and multi-tenancy.

Use this scope when a later task suggests extra infrastructure. If a feature is in `Later extensions`, keep only the interface seam required by the current project.

## Task 3 — Create the folder structure

### Practice guide

Create this structure. Empty directories do not need placeholder files; create each package when its step starts.

```text
cloudianclaw/
├── cmd/agentkit/          process entry point
├── internal/app/          dependency composition
├── internal/agent/        request, prompt, and loop
├── internal/model/        canonical model/provider types
├── internal/tools/        registry and policy
├── internal/session/      history and summaries
├── internal/memory/       added in Steps 8–9
├── internal/runtime/      queue and domain events
├── internal/transport/
│   ├── http/
│   └── websocket/
├── internal/store/postgres/
├── migrations/
├── skills/                local SKILL.md packages
├── web/                   React application
├── go.mod
└── docker-compose.yml
```

Use `cmd/agentkit` as the only executable. Keep `main.go` small: load configuration, build the app, run the selected command, and report startup errors.

## Task 4 — Initialize a buildable Go module

### Practice guide

Initialize the module and add a minimal entry point:

```go
package main

import "fmt"

func main() {
	fmt.Println("agentkit: setup complete")
}
```

Run:

```bash
go test ./...
go run ./cmd/agentkit
```

Expected output:

```text
agentkit: setup complete
```

## Task 5 — Check the architecture boundary

### Practice guide

Review the tree using these rules:

1. `internal/agent` does not import HTTP, WebSocket, SQL drivers, or OpenAI wire types.
2. `internal/transport` does not contain prompt or tool-selection logic.
3. `internal/store/postgres` contains SQL details; consumer packages expose small interfaces.
4. Socket connections, API secrets, and raw media are not modeled as durable messages.
5. Run-specific messages, counters, cancellation, and spans will belong to each request, not a shared singleton.

This step is complete when the module builds, the folder responsibilities are documented, and the project scope is explicit.
