# Step 17 — Observability and Verification

**Knowledge depth: 6/10**  

Read [10 — Tracing & Observability](../10-tracing-observability.md) before tuning or extending the system. It explains the trace, run-timeline, and metric concepts that make an agent's behavior understandable rather than mysterious.

## What GoClaw adds beyond the mini-agent

| Extension | Repository area | Add it when |
|---|---|---|
| Provider registry/fallback | `internal/providerresolve`, `internal/providers` | one adapter and its tests are stable |
| PostgreSQL + SQLite editions | `internal/store/{pg,sqlitestore}` | offline desktop is a real product need |
| MCP bridge/server | `internal/mcp` | tool policy and OAuth lifecycle are mature |
| Skills | `internal/skills` | you need packaged instructions/assets with access control |
| Channels | `internal/channels` | HTTP/WS delivery is reliable |
| Cron/heartbeat | `internal/cron`, `internal/heartbeat` | scheduler drain/cancellation is proven |
| Teams/delegation | `internal/tools/delegate_tool.go`, orchestration packages | task ownership and artifact boundaries exist |
| Knowledge vault/graph | `internal/vault`, `internal/knowledgegraph` | memory source quality needs document/navigation features |
| Media/browser/sandbox | `pkg/browser`, `internal/sandbox`, `internal/media` | per-tenant workspace isolation is tested |

The framework’s broad feature set is useful as a map, but a clone should earn each feature through a concrete product requirement.

## Observability model

Trace one run across all boundaries:

```text
request_id / run_id
  ├── queue.wait
  ├── agent.context
  ├── llm.call (provider, model, tokens, cache, retry)
  ├── tool.call (name, capability, decision, duration)
  ├── memory.retrieve (count, scores, injected tokens)
  └── session.checkpoint
```

GoClaw records LLM tracing in `internal/tracing`, carries run timeline items, and has optional OpenTelemetry export. Start with structured `slog` fields and a durable run ID; add a tracing backend after you can state the questions it must answer.

## Verification ladder

```text
P0  unit/invariant tests: scope, policy, ordering, parsing
P1  contract tests: HTTP and WebSocket response shapes
P2  integration tests: Postgres, vector retrieval, provider fixtures
P3  scenario tests: user journeys across interfaces
```

The repository follows a similar layered approach under `tests/invariants`, `tests/contracts`, `tests/integration`, and `tests/scenarios`. It also has targeted unit tests beside packages; use both styles.

Observability is the feedback loop for the entire course. It connects a user-visible answer to the model calls, tool calls, retrieval decisions, queue wait, and persistence work that produced it.
