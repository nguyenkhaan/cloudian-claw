# Step 19 — Session Scheduler and Shutdown

**Knowledge depth: 8/10**

Serialize work per session while limiting global concurrency and supporting bounded shutdown.

## Step outcome

Same-session jobs run FIFO, unrelated sessions run concurrently, cancellation propagates, and shutdown completes by deadline.

## Task 1 — Define scheduler guarantees

### Theory

- **Architectural role:** Keyed FIFO queues protect consistency within each session, while a global lane protects process and provider resources. The mechanisms are independent so session A does not block session B.
- **Why:** A global worker pool alone does not serialize one session, and per-session mutexes alone do not limit total load. These guarantees are a public contract that tells callers when work has been accepted.
- **GoClaw reference:** Inspect the scheduler entry point in [`goclaw/internal/scheduler/scheduler.go`](../../goclaw/internal/scheduler/scheduler.go), per-session queues in [`goclaw/internal/scheduler/queue.go`](../../goclaw/internal/scheduler/queue.go), and the lane manager in [`goclaw/internal/scheduler/lanes.go`](../../goclaw/internal/scheduler/lanes.go).

Read [08 — Scheduling & Cron](../08-scheduling-cron.md). Two turns for the same session must not load the same old history and then write conflicting results. A slow session must not block unrelated sessions.

The required design is:

```text
per-session FIFO queues → global concurrency semaphore → agent runner
```

Read [22 — Heartbeat System](../22-heartbeat-system.md) only to understand future scheduled autonomy. This Step does not implement cron or heartbeat.

### Goal

Define ordering, concurrency, backpressure, and shutdown behavior before choosing a goroutine and channel structure.

### Guide to implement

Document these guarantees:

- Jobs with the same `(agent_id, user_id, session_key)` run one at a time in submission order.
- Jobs with different keys may run concurrently.
- The process has one configurable maximum number of active chat runs.
- A full queue returns backpressure instead of blocking forever.
- Shutdown rejects new work and drains or cancels existing work by a deadline.

## Task 2 — Implement the keyed scheduler

### Theory

- **Architectural role:** A map and mutex manage queue lifecycle, bounded channels create backpressure, a semaphore or lane grants execution capacity, and a root context controls process lifecycle.
- **Why:** Acquiring a lane only when a job reaches the front of its queue prevents an idle or waiting queue from holding capacity. Cleaning up empty queues prevents the map from growing forever with session keys.
- **GoClaw reference:** Follow `Scheduler.ScheduleWithOpts` and `getOrCreateSession` in [`goclaw/internal/scheduler/scheduler.go`](../../goclaw/internal/scheduler/scheduler.go), then read queue processing in [`goclaw/internal/scheduler/queue.go`](../../goclaw/internal/scheduler/queue.go).

### Goal

Implement actor-like serialization: every session key has a FIFO mailbox, while execution uses a global lane.

### Guide to implement

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

- **Architectural role:** Each queue item owns a request context and one-shot result channel. Its execution context is cancelled by either the request or scheduler shutdown.
- **Why:** Running an already-cancelled job wastes tokens and may cause unwanted tool side effects. Sending exactly one result and then closing the channel gives callers a simple protocol without endless waits or duplicate terminal outcomes.
- **GoClaw reference:** Inspect queue cancellation and result delivery in [`goclaw/internal/scheduler/queue.go`](../../goclaw/internal/scheduler/queue.go), public cancellation in [`goclaw/internal/scheduler/scheduler.go`](../../goclaw/internal/scheduler/scheduler.go), and tests in [`goclaw/internal/scheduler/scheduler_queue_test.go`](../../goclaw/internal/scheduler/scheduler_queue_test.go).

The queued job has two cancellation sources: the request and process shutdown. Cancellation while waiting must not consume a lane or run the model later.

### Goal

Keep a queued job's lifecycle connected to both its caller and the process, even before execution begins.

### Guide to implement

For each scheduled item:

1. Store its request context and result channel.
2. Before execution, skip it if the request is already cancelled.
3. Derive a run context cancelled by either request or shutdown.
4. Send exactly one `JobResult` and close the result channel.
5. Record queue wait time and run time.

Do not discard job errors inside a queue goroutine.

## Task 4 — Implement graceful shutdown

### Theory

- **Architectural role:** Stop the edge first, drain or cancel the scheduler next, drain the background bus after synchronous runs, and close resource clients last.
- **Why:** Closing the database or provider before active runs finish creates artificial failures. Continuing to accept work while draining can prevent shutdown from ever completing. A deadline makes shutdown bounded.
- **GoClaw reference:** Read process lifecycle code in [`goclaw/cmd/gateway_lifecycle.go`](../../goclaw/cmd/gateway_lifecycle.go), scheduler `MarkDraining` and `Stop` in [`goclaw/internal/scheduler/scheduler.go`](../../goclaw/internal/scheduler/scheduler.go), and the gateway `Start` lifecycle in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go).

### Goal

Shut down in dependency order so the system does not accept new work while lower-level dependencies are being removed.

### Guide to implement

Shutdown in this order:

1. Stop accepting HTTP and WebSocket work.
2. Mark the scheduler as draining.
3. Wait for active runs until the shutdown deadline.
4. Cancel remaining runs.
5. Drain background events within a second deadline.
6. Close database and provider clients.

## Task 5 — Verify ordering and concurrency

### Theory

- **Architectural role:** Test jobs use channels or barriers to expose start and finish events, an atomic counter measures global activity, and cancellation and shutdown cover terminal paths.
- **Why:** Concurrency tests based on `time.Sleep` are flaky and do not prove ordering. Controlled jobs let the test pause at the exact transition being asserted.
- **GoClaw reference:** Review [`goclaw/internal/scheduler/scheduler_test.go`](../../goclaw/internal/scheduler/scheduler_test.go) and [`goclaw/internal/scheduler/scheduler_queue_test.go`](../../goclaw/internal/scheduler/scheduler_queue_test.go) for queue and lane testing patterns.

### Goal

Prove scheduler guarantees with controlled synchronization instead of unreliable timing and sleeps.

### Guide to implement

Use controlled jobs with channels to prove:

- Two jobs for session A start and finish in FIFO order.
- A job for session B runs while session A is active when a lane is free.
- Global active jobs never exceed the configured limit.
- A cancelled queued job never starts.
- A full queue returns a clear error.
- Shutdown rejects new jobs and completes within its deadline.

This step is complete when concurrent requests cannot corrupt one session's history.
