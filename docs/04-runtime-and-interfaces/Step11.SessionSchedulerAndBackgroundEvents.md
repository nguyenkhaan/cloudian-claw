# Step 11 — Session Scheduler and Background Events

**Knowledge depth: 7/10**

This step serializes turns in one session, limits global concurrency, and moves memory consolidation out of the request path.

## Task 1 — Define scheduler guarantees

### Theory

Read [08 — Scheduling & Cron](../08-scheduling-cron.md). Two turns for the same session must not load the same old history and then write conflicting results. A slow session must not block unrelated sessions.

The required design is:

```text
per-session FIFO queues → global concurrency semaphore → agent runner
```

Read [22 — Heartbeat System](../22-heartbeat-system.md) only to understand future scheduled autonomy. This Step does not implement cron or heartbeat.

### Practice guide

Document these guarantees:

- Jobs with the same `(agent_id, user_id, session_key)` run one at a time in submission order.
- Jobs with different keys may run concurrently.
- The process has one configurable maximum number of active chat runs.
- A full queue returns backpressure instead of blocking forever.
- Shutdown rejects new work and drains or cancels existing work by a deadline.

## Task 2 — Implement the keyed scheduler

### Practice guide

Create `internal/runtime/scheduler.go`:

```go
type Job func(context.Context) (agent.RunResult, error)

type JobResult struct {
	Value agent.RunResult
	Err   error
}

type Scheduler interface {
	Schedule(ctx context.Context, key string, job Job) <-chan JobResult
	Shutdown(ctx context.Context) error
}
```

The concrete scheduler needs:

- A mutex-protected map of session queues.
- A bounded channel for each queue.
- A semaphore for global concurrency.
- A root context and draining flag.
- Cleanup for inactive empty queues.

Acquire the global lane only when a queued job is ready to run. Always release it with `defer`.

## Task 3 — Propagate cancellation and results

### Theory

The queued job has two cancellation sources: the request and process shutdown. Cancellation while waiting must not consume a lane or run the model later.

### Practice guide

For each scheduled item:

1. Store its request context and result channel.
2. Before execution, skip it if the request is already cancelled.
3. Derive a run context cancelled by either request or shutdown.
4. Send exactly one `JobResult` and close the result channel.
5. Record queue wait time and run time.

Do not discard job errors inside a queue goroutine.

## Task 4 — Implement a small domain event bus

### Theory

Read asynchronous visibility guidance in [10 — Tracing & Observability](../10-tracing-observability.md). Domain events are facts such as “session completed”; they are not chat messages.

### Practice guide

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

Use a bounded queue and worker pool. Register one `session.completed` subscriber that starts Step 9 memory consolidation. Use `SourceID` for handler idempotency.

## Task 5 — Choose delivery semantics

### Theory

Do not claim exactly-once delivery unless both storage and receiver support it.

### Practice guide

Use these semantics:

| Work | Semantic | Recovery |
|---|---|---|
| Chat job | Ordered, at most once after acceptance. | User retries with an idempotency key. |
| Memory consolidation | At least once. | Stable source ID and database UPSERT. |
| WebSocket notification | Best effort. | Client reloads durable state. |

If losing consolidation events during a crash becomes unacceptable, add a transactional outbox as a later extension.

## Task 6 — Implement graceful shutdown

### Practice guide

Shutdown in this order:

1. Stop accepting HTTP and WebSocket work.
2. Mark the scheduler as draining.
3. Wait for active runs until the shutdown deadline.
4. Cancel remaining runs.
5. Drain background events within a second deadline.
6. Close database and provider clients.

## Task 7 — Verify ordering and concurrency

### Practice guide

Use controlled jobs with channels to prove:

- Two jobs for session A start and finish in FIFO order.
- A job for session B runs while session A is active when a lane is free.
- Global active jobs never exceed the configured limit.
- A cancelled queued job never starts.
- A full queue returns a clear error.
- Shutdown rejects new jobs and completes within its deadline.

This step is complete when concurrent requests cannot corrupt one session's history.
