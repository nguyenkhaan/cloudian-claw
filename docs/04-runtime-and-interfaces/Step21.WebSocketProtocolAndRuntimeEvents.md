# Step 21 — WebSocket Protocol and Runtime Events

**Knowledge depth: 8/10**

Add a versioned real-time protocol for chat, tools, and run lifecycle events.

## Step outcome

Clients can start a run, observe ordered progress, disconnect, and recover from durable state.

## Task 1 — Define an explicit socket protocol

### Theory

- **Architectural role:** `req` and `res` form an RPC pair with a request ID, while `event` is server push and does not pretend to be a response. Versions and error codes create a contract that can evolve and be tested separately from transport implementation.
- **Why:** Ad hoc JSON forces clients to guess shapes and can attach responses or events incorrectly. An explicit discriminator enables early rejection and deliberate backward compatibility.
- **GoClaw reference:** Compare `RawFrame`, `RequestFrame`, `ResponseFrame`, and `EventFrame` in [`goclaw/pkg/protocol/frames.go`](../../goclaw/pkg/protocol/frames.go), method constants in [`goclaw/pkg/protocol/methods.go`](../../goclaw/pkg/protocol/methods.go), and the router in [`goclaw/internal/gateway/router.go`](../../goclaw/internal/gateway/router.go).

Read [19 — WebSocket RPC Methods](../19-websocket-rpc.md) and [04 — Gateway and Protocol](../04-gateway-protocol.md). [13 — WebSocket Team & Delegation Events](../13-ws-team-events.md) is reference material; team events are not implemented.

Separate request, response, and pushed event frames:

```json
{"type":"req","id":"r-42","method":"chat.send","params":{}}
{"type":"res","id":"r-42","ok":true,"payload":{}}
{"type":"event","event":"chat.delta","payload":{"text":"Hello"},"seq":12}
```

### Goal

Define frame envelopes and correlation rules before writing connection code or a UI client.

### Guide to implement

Create shared protocol types with a version field. Reject unknown frame types, missing request IDs, unknown methods, and oversized frames with stable error codes.

## Task 2 — Define the event catalogue

### Theory

- **Architectural role:** Run, session, and call IDs are correlation keys, `seq` defines order, and a terminal event closes the lifecycle. Payloads contain observable facts rather than internal objects or pointers.
- **Why:** Text deltas alone cannot distinguish two runs or describe tool activity. A typed catalogue prevents each producer from inventing its own event names and shapes.
- **GoClaw reference:** Inspect event constructors and types in [`goclaw/pkg/protocol/events.go`](../../goclaw/pkg/protocol/events.go), run timeline events in [`goclaw/internal/agent/run_timeline_recorder.go`](../../goclaw/internal/agent/run_timeline_recorder.go), and UI event dispatch in [`goclaw/ui/web/src/hooks/use-ws-event.ts`](../../goclaw/ui/web/src/hooks/use-ws-event.ts).

### Goal

Turn run progress into a stable vocabulary understood by both the backend and UI.

### Guide to implement

Start with these events:

| Event | Required payload |
|---|---|
| `run.started` | `run_id`, `session_key` |
| `chat.delta` | `run_id`, `session_key`, `text`, `seq` |
| `tool.started` | `run_id`, `tool_call_id`, `name` |
| `tool.completed` | `run_id`, `tool_call_id`, `name`, `is_error` |
| `run.completed` | `run_id`, `message_id`, `usage` |
| `run.failed` | `run_id`, safe `code`, safe `message` |

Sequence numbers are monotonic per connection or run; document which choice you use. Include run and session IDs so the UI cannot attach late events to the wrong view.

## Task 3 — Implement connection lifecycle

### Theory

- **Architectural role:** One reader parses requests, one writer owns every write, a bounded queue separates producers from the writer, and the connection context owns subscriptions and is cancelled on close.
- **Why:** Concurrent writes can corrupt protocol frames, an unbounded queue lets a slow client consume all memory, and ping or idle deadlines clean up dead connections that TCP has not yet reported.
- **GoClaw reference:** Inspect upgrade and server lifecycle in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go), per-client reading and writing in [`goclaw/internal/gateway/client.go`](../../goclaw/internal/gateway/client.go), and origin and rate-limit tests in [`goclaw/internal/gateway/server_test.go`](../../goclaw/internal/gateway/server_test.go).

One reader and one writer avoid concurrent socket writes and simplify backpressure. The first successful action must establish trusted identity.

### Goal

Manage authentication, concurrency, liveness, and backpressure for each socket as a resource with a lifecycle.

### Guide to implement

For each connection:

1. Authenticate the HTTP upgrade request.
2. Check the allowed origin.
3. Set frame and read limits.
4. Start one read loop and one bounded write queue/loop.
5. Add ping/pong and idle deadlines.
6. Cancel connection-owned subscriptions on close.

When the write queue is full, close the slow connection with a clear reason. Do not block agent execution forever.

## Task 4 — Implement `chat.send`

### Theory

- **Architectural role:** The method handler validates and submits, callbacks or an event publisher report progress, and durable storage determines completion. The socket observes the run but does not own it.
- **Why:** Keeping the RPC open until the run finishes makes correlation, cancellation, and reconnection harder. Emitting completion before persistence creates a false completion if the process stops immediately afterward.
- **GoClaw reference:** Follow the chat method in [`goclaw/internal/gateway/methods/chat.go`](../../goclaw/internal/gateway/methods/chat.go), chat behavior helpers in [`goclaw/internal/gateway/methods/chat_behavior.go`](../../goclaw/internal/gateway/methods/chat_behavior.go), and event broadcasting in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go).

### Goal

Start an asynchronous run through the same use case as HTTP and return an acknowledgement immediately.

### Guide to implement

Validate params, call the same `ChatService` as HTTP, and return an immediate response containing `run_id`. Stream later progress as events.

Connect Step 7 streaming callbacks and Step 13 tool/run callbacks to the event publisher. Persist durable messages before emitting `run.completed`.

Do not write socket references into session or run tables.

## Task 5 — Support reconnection

### Theory

- **Architectural role:** WebSocket optimizes latency, HTTP plus session and run stores provide authoritative snapshots, and the client reconciles the snapshot before consuming new events.
- **Why:** Treating an in-memory event stream as a durable log requires complex replay offsets and storage. Reloading durable state is simpler and correct for the first version.
- **GoClaw reference:** Inspect UI reconnection in [`goclaw/ui/web/src/api/ws-client.ts`](../../goclaw/ui/web/src/api/ws-client.ts), the provider wrapper in [`goclaw/ui/web/src/components/providers/ws-provider.tsx`](../../goclaw/ui/web/src/components/providers/ws-provider.tsx), and durable session queries in [`goclaw/internal/gateway/methods/sessions.go`](../../goclaw/internal/gateway/methods/sessions.go).

Events are a live view, not the database. A client can disconnect between deltas.

### Goal

Recover the UI from durable truth after best-effort live events are lost.

### Guide to implement

The UI recovery flow is:

```text
reconnect socket
→ authenticate
→ reload session messages and run state through HTTP
→ subscribe to new events
```

The first version does not need event replay. Durable HTTP queries are the recovery path.

## Task 6 — Verify protocol behavior

### Theory

- **Architectural role:** Protocol tests assert raw frames, server tests assert authentication and limits, and a reconnect scenario verifies that the event view converges on durable state.
- **Why:** Component tests may not catch a wrong response ID or a terminal event arriving before its deltas. A slow-client test proves that the bounded queue really protects the server.
- **GoClaw reference:** Review [`goclaw/internal/gateway/server_test.go`](../../goclaw/internal/gateway/server_test.go), [`goclaw/internal/gateway/router_test.go`](../../goclaw/internal/gateway/router_test.go), the [`goclaw/pkg/protocol`](../../goclaw/pkg/protocol) package, and the UI WebSocket client in [`goclaw/ui/web/src/api/ws-client.ts`](../../goclaw/ui/web/src/api/ws-client.ts).

### Goal

Test contract, ordering, correlation, pressure, and recovery across the complete connection lifecycle.

### Guide to implement

Test:

1. Authenticated and rejected upgrades.
2. Request ID matched to response ID.
3. Ordered text deltas and terminal event.
4. Tool events correlated by call ID.
5. Oversized/invalid frame rejection.
6. Slow-client backpressure.
7. Reconnect followed by successful history reload.

This step is complete when a client can start one run, observe it live, disconnect, and recover authoritative state.
