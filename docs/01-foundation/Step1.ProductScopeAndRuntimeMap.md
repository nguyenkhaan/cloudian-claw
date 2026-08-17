# Step 1 — Product Scope and Runtime Map

**Knowledge depth: 4/10**

Define the product boundary and understand how a request moves through the system.

## Step outcome

The project has an explicit first-release scope and a package owner for every runtime responsibility.

## Task 1 — Understand the runtime path

### Theory

- **Architectural role:** This is the high-level dependency map. The transport receives and returns data, the scheduler coordinates work, the agent makes decisions, providers and tools perform external work, and stores keep durable state.
- **Why this comes first:** Without this map, agent logic can easily end up in HTTP handlers, SQL can be called directly from the agent loop, and the HTTP, CLI, and WebSocket paths can behave differently.
- **GoClaw reference:** Start at [`goclaw/cmd/gateway.go`](../../goclaw/cmd/gateway.go), follow dependency wiring in [`goclaw/cmd/gateway_setup.go`](../../goclaw/cmd/gateway_setup.go), request handling in [`goclaw/internal/gateway/chat_runner.go`](../../goclaw/internal/gateway/chat_runner.go), and then [`Loop.Run`](../../goclaw/internal/agent/loop_run.go). GoClaw is larger than this course project, but the chain of responsibilities is similar.

Read [00 — Architecture Overview](../00-architecture-overview.md). Focus on the startup path, package responsibilities, and the boundary between the agent core and GoClaw's platform features.

GoClaw is a gateway, not only an LLM wrapper. A request enters through a transport, waits in a session queue, runs through an agent pipeline, may call tools, persists durable state, and returns through the same transport.

```text
request → session queue → agent pipeline → provider/tools → store → response
```

### Goal

Draw the path of a request and identify the package responsible for each stage before writing code.

### Guide to implement

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

- **Architectural role:** This scope keeps the important extension seams (`Provider`, `Store`, `Tool`, and `ChatService`) without bringing all production complexity into the first version.
- **Why:** GoClaw supports multi-tenancy, many channels and providers, teams, and cron jobs. Building all of them at once would hide core invariants such as message order, user scope, and tool authorization.
- **GoClaw reference:** Use the project map in [`goclaw/AGENTS.md`](../../goclaw/AGENTS.md) to see the full production scope, then compare the core with [`goclaw/internal/agent`](../../goclaw/internal/agent) and [`goclaw/internal/providers`](../../goclaw/internal/providers). The `Later extensions` list represents GoClaw packages that this course intentionally postpones.

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

### Goal

Choose one vertical slice that is deep enough to teach the complete flow instead of partially building many unrelated features.

### Guide to implement

Add a short `Scope` section to the project README with two lists:

- `Build now`: the seven items above.
- `Later extensions`: channels, teams, provider routing, skill publishing, autonomous schedules, and multi-tenancy.

Use this scope when a later task suggests extra infrastructure. If a feature is in `Later extensions`, keep only the interface seam required by the current project.
