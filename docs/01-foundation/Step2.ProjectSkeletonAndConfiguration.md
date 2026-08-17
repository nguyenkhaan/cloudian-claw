# Step 2 — Project Skeleton and Configuration

**Knowledge depth: 4/10**

Create the buildable project skeleton and one validated configuration entry point.

## Step outcome

The Go module builds, the composition boundaries are visible, and invalid configuration fails at startup.

## Task 1 — Create the folder structure

### Theory

- **Architectural role:** The structure follows a ports-and-adapters style: `agent/model` is the core, `providers/store/tools` are adapters, `transport` contains inbound adapters, and `app` is the composition root.
- **Why:** Packages do more than organize files. Import direction determines whether the core can be tested with fakes and whether infrastructure can be replaced without changing business behavior.
- **GoClaw reference:** Compare [`goclaw/internal/agent`](../../goclaw/internal/agent), [`goclaw/internal/store`](../../goclaw/internal/store), [`goclaw/internal/gateway`](../../goclaw/internal/gateway), [`goclaw/pkg/protocol`](../../goclaw/pkg/protocol), and composition in [`goclaw/cmd/gateway_setup.go`](../../goclaw/cmd/gateway_setup.go). CloudianClaw uses fewer packages but keeps the same dependency direction.

### Goal

Turn the logical boundaries from Task 1 into package boundaries that the compiler can help enforce.

### Guide to implement

Create this structure. Empty directories do not need placeholder files; create each package when its step starts.

```text
cloudianclaw/
├── cmd/agentkit/          process entry point
├── internal/app/          dependency composition
├── internal/agent/        request, prompt, and loop
├── internal/model/        canonical model/provider types
├── internal/tools/        registry and policy
├── internal/session/      history and summaries
├── internal/memory/       added in Steps 16–18
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

## Task 2 — Initialize a buildable Go module

### Theory

- **Architectural role:** `cmd/agentkit` is the process entry point, not a home for business logic. Later, it should only parse commands, call the composition root, and manage process lifecycle.
- **Why:** A thin executable lets packages be tested without starting the full server and avoids multiple, inconsistent ways to construct dependencies.
- **GoClaw reference:** Read [`goclaw/main.go`](../../goclaw/main.go) and [`goclaw/cmd/root.go`](../../goclaw/cmd/root.go) to see how the entry point delegates to the CLI. Most dependency wiring is split across the `goclaw/cmd/gateway_*.go` files.

### Goal

Create a buildable and testable baseline so every later Task adds one small, verifiable change.

### Guide to implement

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

## Task 3 — Check the architecture boundary

### Theory

- **Architectural role:** This is an architecture fitness check: the core knows only contracts, adapters know external details, and request-specific state never lives on a shared singleton.
- **Why:** Boundary mistakes often compile successfully and appear only after adding another provider or transport, or under concurrency. Checking imports and ownership early is cheaper than a later refactor.
- **GoClaw reference:** Compare the provider contract in [`goclaw/internal/providers/types.go`](../../goclaw/internal/providers/types.go), the store interface in [`goclaw/internal/store/session_store.go`](../../goclaw/internal/store/session_store.go), and per-run state in [`goclaw/internal/pipeline/run_state.go`](../../goclaw/internal/pipeline/run_state.go). Notice that GoClaw also keeps WebSocket connections in the gateway rather than in session data.

### Goal

Turn package decisions into a reviewable checklist instead of assuming that the folder structure looks reasonable.

### Guide to implement

Review the tree using these rules:

1. `internal/agent` does not import HTTP, WebSocket, SQL drivers, or OpenAI wire types.
2. `internal/transport` does not contain prompt or tool-selection logic.
3. `internal/store/postgres` contains SQL details; consumer packages expose small interfaces.
4. Socket connections, API secrets, and raw media are not modeled as durable messages.
5. Run-specific messages, counters, cancellation, and spans will belong to each request, not a shared singleton.

This step is complete when the module builds, the folder responsibilities are documented, and the project scope is explicit.

## Task 4 — Load and validate configuration

### Theory

- **Architectural role:** Configuration belongs at the composition root. Domain packages receive only the small config values they need instead of reading environment variables themselves. Secrets are runtime inputs, not domain data.
- **Why:** Failing fast during startup is better than discovering a missing URL or model halfway through a request. Config groups also clarify which settings belong to HTTP, providers, agents, or the runtime.
- **GoClaw reference:** Inspect the grouped structure in [`goclaw/internal/config/config.go`](../../goclaw/internal/config/config.go), defaults, loading, and environment overlays in [`goclaw/internal/config/config_load.go`](../../goclaw/internal/config/config_load.go), and secret-specific configuration in [`goclaw/internal/config/config_secrets.go`](../../goclaw/internal/config/config_secrets.go).

### Goal

Collect runtime configuration into one validated object before constructing dependencies.

### Guide to implement

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
