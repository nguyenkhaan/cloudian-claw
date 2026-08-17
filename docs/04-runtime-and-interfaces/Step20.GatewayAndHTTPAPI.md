# Step 20 — Gateway and HTTP API

**Knowledge depth: 7/10**

Expose the application service through one authenticated and validated HTTP boundary.

## Step outcome

Health and chat endpoints use trusted identity, scheduler-backed execution, structured logs, and stable errors.

## Task 1 — Define the gateway boundary

### Theory

- **Architectural role:** HTTP, CLI, and WebSocket are adapters, `ChatService` is the inbound port, and the scheduler plus runner form the application core. Trusted identity becomes a scoped `RunRequest` at this boundary.
- **Why:** If every transport constructs its own runner, key, and policy, persistence and authorization behavior will drift. One shared service keeps transports thin and the contract testable.
- **GoClaw reference:** Inspect the gateway chat bridge in [`goclaw/internal/gateway/chat_runner.go`](../../goclaw/internal/gateway/chat_runner.go), the WebSocket method in [`goclaw/internal/gateway/methods/chat.go`](../../goclaw/internal/gateway/methods/chat.go), and the HTTP adapter in [`goclaw/internal/http/chat_completions.go`](../../goclaw/internal/http/chat_completions.go).

Read [04 — Gateway and Protocol](../04-gateway-protocol.md) and [18 — HTTP REST API](../18-http-api.md). A transport authenticates, parses, validates, calls the application service, and serializes a result. It does not contain agent behavior.

```text
CLI/HTTP → chat service → scheduler → agent runner → stores
```

### Goal

Create one application use case called by every inbound adapter.

### Guide to implement

Create an application service used by both transports:

```go
type ChatService interface {
	Submit(context.Context, agent.RunRequest) <-chan runtime.JobResult
}
```

The service builds the scheduler key from trusted agent/user identity and the validated session key.

## Task 2 — Add edge authentication

### Theory

- **Architectural role:** Middleware authenticates and stores the principal and request ID in context. Downstream code authorizes resources with that principal instead of reading the token again or trusting `user_id` from the payload.
- **Why:** Scattered identity resolution creates bypasses and inconsistent errors. Constant-time comparison reduces timing side channels, and a generic 401 response does not reveal information about the token.
- **GoClaw reference:** Study authentication helpers and middleware in [`goclaw/internal/http/auth.go`](../../goclaw/internal/http/auth.go), gateway connection authentication in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go), and actor and tenant context helpers in [`goclaw/internal/store/context.go`](../../goclaw/internal/store/context.go).

Use [20 — API Keys & Authentication](../20-api-keys-auth.md) for the larger design. This project uses one configured bearer token and one server-side user identity.

### Goal

Convert an untrusted credential into a trusted principal exactly once at the edge.

### Guide to implement

Implement middleware that:

1. Reads `Authorization: Bearer <token>`.
2. Compares it in constant time with the configured token.
3. Creates a trusted principal with user and allowed agent IDs.
4. Adds a request ID and principal to context.
5. Returns `401` without revealing why a token failed.

Do not accept `user_id` from request JSON as authority.

## Task 3 — Implement health and chat endpoints

### Theory

- **Architectural role:** The handler performs parse → validate → principal and resource resolution → service call → serialization. It does not build prompts, call providers, or write sessions directly.
- **Why:** Body, field, and length limits stop resource abuse before it reaches the core. An error taxonomy separates client, queue, upstream, and timeout failures so retry and UX behavior can be correct without exposing internals.
- **GoClaw reference:** Inspect request, response, and streaming paths in [`goclaw/internal/http/chat_completions.go`](../../goclaw/internal/http/chat_completions.go), response helpers in [`goclaw/internal/http/responses.go`](../../goclaw/internal/http/responses.go), health wiring in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go), and contract tests in [`goclaw/internal/http/api_contracts_test.go`](../../goclaw/internal/http/api_contracts_test.go).

### Goal

Expose the use case through HTTP with input bounds, typed validation, and stable error mapping.

### Guide to implement

Add:

```text
GET  /health
POST /v1/chat/completions
```

Use a small request shape:

```json
{
  "agent_id": "agent-uuid",
  "session_key": "demo",
  "message": "Summarize our last decision."
}
```

The handler must:

- Limit body size.
- Reject unknown JSON fields.
- Validate message and session-key length.
- Resolve user identity from context.
- Wait for the scheduled result or request cancellation.
- Return JSON with content, usage, iteration count, and tool-call count.

Map validation to `400`, authentication to `401`, missing resources to `404`, queue pressure to `429`, and upstream failure to a safe `502`/`504` response.

## Task 4 — Add structured request logs

### Theory

- **Architectural role:** The edge creates `request_id`, the application creates `run_id`, and completion logs record duration, status, and error code. IDs travel through context and events to connect records.
- **Why:** Logging prompts and tool output consumes space and creates data leaks. Structured metadata supports filtering and metrics better than arbitrary log strings.
- **GoClaw reference:** Inspect logging and error helpers around [`goclaw/internal/gateway/log_tee.go`](../../goclaw/internal/gateway/log_tee.go), tracing integration in [`goclaw/internal/agent/loop_tracing.go`](../../goclaw/internal/agent/loop_tracing.go), and HTTP response and error handling under [`goclaw/internal/http`](../../goclaw/internal/http).

### Goal

Create enough correlation data to debug a request across transport, queue, and run without logging sensitive content.

### Guide to implement

Log one start and one completion record containing:

```text
request_id, method, path, user_id, agent_id, session_key,
status, duration_ms, run_id, error_code
```

Do not log bearer tokens, provider keys, complete prompts, or full tool results.
