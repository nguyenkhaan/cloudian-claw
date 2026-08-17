# Step 13 — WebSocket Events and UI Contract

**Knowledge depth: 8/10**

Read [19 — WebSocket RPC Methods](../19-websocket-rpc.md) and [04 — Gateway and Protocol](../04-gateway-protocol.md) before introducing a socket. They explain GoClaw's request, response, and event frame shapes. [13 — WebSocket Team & Delegation Events](../13-ws-team-events.md) is useful architecture reading, but this course emits only chat text, tool activity, and run-status events.

## Why a chat UI needs more than HTTP

HTTP is excellent for a completed result. A user interface also benefits from token deltas, tool activity, status changes, and background task updates. A WebSocket provides a long-lived, bidirectional channel for those events.

```text
client request frame → gateway method → scheduler → agent run
agent/runner events → gateway event frame → client state update
```

## Keep the protocol explicit

GoClaw uses three frame categories:

```json
{"type":"req","id":"r-42","method":"connect","params":{}}
{"type":"res","id":"r-42","ok":true,"payload":{}}
{"type":"event","event":"chat.delta","payload":{"text":"Hello"},"seq":12}
```

The client-generated request ID matches a response. Events have their own name and may carry a sequence number. This small distinction avoids confusing a pushed update with a response to an action.

## Connection lifecycle

The first method establishes identity and connection scope. After that, each RPC method follows the same path as the HTTP API: validate input, call the shared application service, and serialize a typed result. One read loop and one write loop keep socket I/O orderly.

## Events are a view, not the database

A client can disconnect between two events. Design the UI so it can re-query durable state—sessions, messages, tasks, or traces—after reconnecting. The event stream makes interaction immediate; the store remains authoritative.
