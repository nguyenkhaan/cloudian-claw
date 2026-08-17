# Step 13 — WebSocket Events and UI Contract

**Knowledge depth: 8/10**

This step adds real-time text, tool activity, and run-status events for the web UI.

## Task 1 — Define an explicit socket protocol

### Theory

Read [19 — WebSocket RPC Methods](../19-websocket-rpc.md) and [04 — Gateway and Protocol](../04-gateway-protocol.md). [13 — WebSocket Team & Delegation Events](../13-ws-team-events.md) is reference material; team events are not implemented.

Separate request, response, and pushed event frames:

```json
{"type":"req","id":"r-42","method":"chat.send","params":{}}
{"type":"res","id":"r-42","ok":true,"payload":{}}
{"type":"event","event":"chat.delta","payload":{"text":"Hello"},"seq":12}
```

### Practice guide

Create shared protocol types with a version field. Reject unknown frame types, missing request IDs, unknown methods, and oversized frames with stable error codes.

## Task 2 — Define the event catalogue

### Practice guide

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

One reader and one writer avoid concurrent socket writes and simplify backpressure. The first successful action must establish trusted identity.

### Practice guide

For each connection:

1. Authenticate the HTTP upgrade request.
2. Check the allowed origin.
3. Set frame and read limits.
4. Start one read loop and one bounded write queue/loop.
5. Add ping/pong and idle deadlines.
6. Cancel connection-owned subscriptions on close.

When the write queue is full, close the slow connection with a clear reason. Do not block agent execution forever.

## Task 4 — Implement `chat.send`

### Practice guide

Validate params, call the same `ChatService` as HTTP, and return an immediate response containing `run_id`. Stream later progress as events.

Connect Step 4 streaming callbacks and Step 6 tool/run callbacks to the event publisher. Persist durable messages before emitting `run.completed`.

Do not write socket references into session or run tables.

## Task 5 — Support reconnection

### Theory

Events are a live view, not the database. A client can disconnect between deltas.

### Practice guide

The UI recovery flow is:

```text
reconnect socket
→ authenticate
→ reload session messages and run state through HTTP
→ subscribe to new events
```

The first version does not need event replay. Durable HTTP queries are the recovery path.

## Task 6 — Verify protocol behavior

### Practice guide

Test:

1. Authenticated and rejected upgrades.
2. Request ID matched to response ID.
3. Ordered text deltas and terminal event.
4. Tool events correlated by call ID.
5. Oversized/invalid frame rejection.
6. Slow-client backpressure.
7. Reconnect followed by successful history reload.

This step is complete when a client can start one run, observe it live, disconnect, and recover authoritative state.
