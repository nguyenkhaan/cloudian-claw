# Step 1 — Folder Structure Initialization

**Knowledge depth: 4/10**  

Start by reading [00 — Architecture Overview](../00-architecture-overview.md). Focus on the startup path, package responsibilities, and the boundary between a small agent core and GoClaw's larger platform.

GoClaw is not just an LLM wrapper. It is a multi-surface agent gateway: a request enters through a channel, CLI, HTTP, or WebSocket; a scheduler protects ordering; an agent pipeline calls a provider and tools; stores keep context and memory; results are delivered back through the original surface.

## Begin with a useful slice

Build this first:

```mermaid
flowchart LR
    I[CLI or HTTP request] --> S[Session queue]
    S --> A[Agent pipeline]
    A --> P[LLM provider]
    P -->|tool calls| T[Tool registry]
    T --> A
    A --> D[(Session store)]
    A --> O[Response]
```

Everything else is an extension of one of those arrows.

| Concern | Minimal version | GoClaw reference |
|---|---|---|
| Entry | HTTP `POST /chat` | `internal/http`, `internal/gateway`, `cmd` |
| Execution | one loop per session | `internal/pipeline`, `internal/agent` |
| Model | one OpenAI-compatible provider | `internal/providers`, `internal/providerresolve` |
| State | messages and a summary | `internal/store/session_store.go` |
| Action | read-only `time` tool | `internal/tools` |
| Persistence | Postgres | `internal/store/pg`, `migrations/` |

## Read the real architecture in this order

```text
main.go
  cmd/root.go, cmd/gateway.go           process composition
  internal/gateway/server.go            HTTP + WebSocket boundary
  internal/agent/router.go              agent lookup/cache
  internal/agent/loop_pipeline_adapter.go
  internal/pipeline/pipeline.go         execution flow
  internal/providers/types.go           model-neutral contract
  internal/tools/types.go               tool contract
  internal/store/session_store.go       durable conversation contract
```

The current runtime is a pipeline. Its public comment names eight logical phases:

```text
context → history → prompt → think → act → observe → memory → summarize
```

In code, some phases are grouped: `ContextStage` prepares context/history/prompt; `ThinkStage` calls the model; `ToolStage` acts; `ObserveStage` records model output; pruning and checkpoint stages manage context and persistence.

## Separate core from platform features

This separation prevents a common failure: reproducing directories instead of reproducing behavior.

```text
Core agent
├── canonical messages
├── provider interface
├── bounded think → tool → observe loop
├── session persistence
└── tool authorization

Platform features
├── multi-tenancy / RBAC / encrypted credentials
├── Telegram, Discord, WhatsApp, Feishu, Zalo, and other channels
├── MCP servers and OAuth
├── browser, sandbox, media, TTS
├── teams, delegation, cron, heartbeat
└── React and desktop clients
```

## Architecture rules to adopt now

1. **Normalize at boundaries.** Providers and tools have vendor-specific shapes; the agent loop must not.
2. **Keep run state per run.** A shared loop object may hold dependencies, but messages, counters, spans, and cancellation must belong to the request.
3. **Serialize a session by default.** Two simultaneous turns on the same history cause incorrect answers and lost writes.
4. **Persist only durable facts.** Do not store socket connections, raw API secrets, or giant base64 media in message history.
5. **Use `context.Context` for scope and cancellation.** GoClaw propagates tenant, agent, user, locale, and workspace through context-aware stores.

As you create the folders, keep the first architecture deliberately small: one process, one provider adapter, a durable session store, and a narrow tool surface. The later steps grow this core into a complete gateway and UI.
