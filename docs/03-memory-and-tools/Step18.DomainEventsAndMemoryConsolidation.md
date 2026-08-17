# Step 18 — Domain Events and Memory Consolidation

**Knowledge depth: 9/10**

Move session-to-memory consolidation out of the request path through explicit domain events.

## Step outcome

Completed sessions produce retry-safe episodic records, and the agent can use both automatic recall and memory_search.

## Task 1 — Implement a small domain event bus

### Theory

- **Architectural role:** A publisher emits a domain fact, the bus manages queues and workers, and a subscriber owns the side effect. Events contain stable identity and provenance, not sockets or behavior callbacks.
- **Why:** Calling consolidation directly increases request coupling and latency. A bounded queue requires an explicit overload policy instead of creating unlimited goroutines.
- **GoClaw reference:** Read the contract and configuration in [`goclaw/internal/eventbus/domain_event_bus.go`](../../goclaw/internal/eventbus/domain_event_bus.go), worker and delivery implementation in [`goclaw/internal/eventbus/bus_impl.go`](../../goclaw/internal/eventbus/bus_impl.go), and the event catalogue in [`goclaw/internal/eventbus/event_types.go`](../../goclaw/internal/eventbus/event_types.go).

Read asynchronous visibility guidance in [10 — Tracing & Observability](../10-tracing-observability.md). Domain events are facts such as “session completed”; they are not chat messages.

### Goal

Move post-processing that does not need to block the response out of the synchronous agent run.

### Guide to implement

Create an in-process bus:

```go
type Event struct {
	Type     string
	SourceID string
	Payload  any
}

type Handler func(context.Context, Event) error

type Bus interface {
	Publish(context.Context, Event) error
	Subscribe(kind string, handler Handler) func()
	Start(context.Context)
	Drain(context.Context) error
}
```

Use a bounded queue and worker pool. Register one `session.completed` subscriber for the consolidation handler implemented later in this Step. Use `SourceID` for handler idempotency.

## Task 2 — Choose delivery semantics

### Theory

- **Architectural role:** Chat ordering belongs to the scheduler, consolidation deduplication belongs to the event handler and store, and notification recovery belongs to the client plus durable queries. End-to-end guarantees are the combination of these layers.
- **Why:** An in-process bus cannot guarantee that events survive a crash. An idempotency key and UPSERT make at-least-once handling practical, while UI events are best effort because the database is the source of truth.
- **GoClaw reference:** Inspect deduplication in [`goclaw/internal/eventbus/dedup.go`](../../goclaw/internal/eventbus/dedup.go), the idempotent episodic handler in [`goclaw/internal/consolidation/episodic_worker.go`](../../goclaw/internal/consolidation/episodic_worker.go), and pushed events in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go).

Do not claim exactly-once delivery unless both storage and receiver support it.

### Goal

Assign an appropriate consistency and delivery guarantee to each kind of work instead of using an “exactly once” slogan.

### Guide to implement

Use these semantics:

| Work | Semantic | Recovery |
|---|---|---|
| Chat job | Ordered, at most once after acceptance. | User retries with an idempotency key. |
| Memory consolidation | At least once. | Stable source ID and database UPSERT. |
| WebSocket notification | Best effort. | Client reloads durable state. |

If losing consolidation events during a crash becomes unacceptable, add a transactional outbox as a later extension.

## Task 3 — Consolidate completed sessions

### Theory

- **Architectural role:** The request path publishes only a domain fact. A worker loads the durable source, extracts, filters, embeds, and upserts it. A stable `SourceID` connects the event to its idempotency key.
- **Why:** Consolidation adds LLM, network, and database work but does not need to delay the response. At-least-once delivery requires an idempotent handler, and UPSERT prevents duplicate records.
- **GoClaw reference:** Follow [`episodicWorker.Handle`](../../goclaw/internal/consolidation/episodic_worker.go), worker registration and expiry cleanup in [`goclaw/internal/consolidation/workers.go`](../../goclaw/internal/consolidation/workers.go), event definitions in [`goclaw/internal/eventbus/event_types.go`](../../goclaw/internal/eventbus/event_types.go), and cleanup contract tests in [`goclaw/internal/consolidation/prune_expired_test.go`](../../goclaw/internal/consolidation/prune_expired_test.go).

Use [10 — Tracing & Observability](../10-tracing-observability.md) to decide which retrieval and consolidation facts to record. Consolidation should not delay the chat response.

### Goal

Turn a session outcome into an asynchronous episodic record that is retry-safe and measurable.

### Guide to implement

When a session is summarized, publish an event with scope, session key, and stable source ID. The handler performs:

```text
load durable summary
→ extract one useful episode
→ validate policy and scope
→ create L0 abstract and topics
→ embed summary
→ UPSERT by stable source ID
→ emit timing/result metrics
```

Start with one handler. Retry transient failures with the same source ID. Periodically delete expired records.

## Task 4 — Connect memory search to the run

### Theory

- **Architectural role:** The automatic path favors low latency and L0 abstracts. The tool path runs during act-observe and returns source metadata. Both use the same scoped store and policy but have different disclosure budgets.
- **Why:** Automatic injection alone can overload context, while tool search alone requires the model to know it needs memory before seeing any hint. The two paths complement each other.
- **GoClaw reference:** Inspect run integration in [`goclaw/internal/agent/loop_context.go`](../../goclaw/internal/agent/loop_context.go), the automatic injector in [`goclaw/internal/memory/auto_injector_impl.go`](../../goclaw/internal/memory/auto_injector_impl.go), and explicit `MemorySearchTool` in [`goclaw/internal/tools/memory.go`](../../goclaw/internal/tools/memory.go).

### Goal

Place automatic recall in the correct pre-prompt phase while also providing an explicit tool for deeper search.

### Guide to implement

Before prompt construction:

1. Build a recall query from the latest message and a short recent frame.
2. Embed it.
3. Run scoped hybrid search.
4. Build the L0 injection for `PromptBuilder`.

Register `memory_search` for deeper, explicit retrieval by the model. Return source IDs and scores for observability, but do not expose another user's identifiers.

## Task 5 — Verify retrieval safety and quality

### Theory

- **Architectural role:** A multi-user fixture checks hard scope, expired and duplicate records check lifecycle, exact and paraphrased queries check lexical and vector search, and injection assertions check downstream budgets.
- **Why:** A single-user relevance benchmark cannot detect leaks, while a scope-only test does not prove that retrieval is useful. Memory provides value only when it is both safe and effective.
- **GoClaw reference:** Review [`goclaw/internal/memory/recall_query_test.go`](../../goclaw/internal/memory/recall_query_test.go), [`goclaw/internal/memory/embeddings_provider_test.go`](../../goclaw/internal/memory/embeddings_provider_test.go), and consolidation tests under [`goclaw/internal/consolidation`](../../goclaw/internal/consolidation).

### Goal

Test security invariants, idempotency, and minimum quality for both retrieval signals together.

### Guide to implement

Seed records for two users and test:

- User A retrieves a relevant paraphrase from A's memory.
- User A never receives user B's record.
- Expired records are excluded.
- Duplicate consolidation events create one record.
- Prompt injection respects item and token limits.
- Exact names can be found through the lexical score.

This step is complete when a new session can recall a scoped fact created from an earlier session.
