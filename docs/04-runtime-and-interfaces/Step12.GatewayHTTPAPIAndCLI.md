# Step 12 — Gateway, HTTP API, and CLI

**Knowledge depth: 7/10**

This step exposes the same agent application service through HTTP and a small CLI.

## Task 1 — Define the gateway boundary

### Theory

Read [04 — Gateway and Protocol](../04-gateway-protocol.md) and [18 — HTTP REST API](../18-http-api.md). A transport authenticates, parses, validates, calls the application service, and serializes a result. It does not contain agent behavior.

```text
CLI/HTTP → chat service → scheduler → agent runner → stores
```

### Practice guide

Create an application service used by both transports:

```go
type ChatService interface {
	Submit(context.Context, agent.RunRequest) <-chan runtime.JobResult
}
```

The service builds the scheduler key from trusted agent/user identity and the validated session key.

## Task 2 — Add edge authentication

### Theory

Use [20 — API Keys & Authentication](../20-api-keys-auth.md) for the larger design. This project uses one configured bearer token and one server-side user identity.

### Practice guide

Implement middleware that:

1. Reads `Authorization: Bearer <token>`.
2. Compares it in constant time with the configured token.
3. Creates a trusted principal with user and allowed agent IDs.
4. Adds a request ID and principal to context.
5. Returns `401` without revealing why a token failed.

Do not accept `user_id` from request JSON as authority.

## Task 3 — Implement health and chat endpoints

### Practice guide

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

### Practice guide

Log one start and one completion record containing:

```text
request_id, method, path, user_id, agent_id, session_key,
status, duration_ms, run_id, error_code
```

Do not log bearer tokens, provider keys, complete prompts, or full tool results.

## Task 5 — Implement CLI commands

### Theory

The CLI is an HTTP client of the running gateway. Building a second in-process runner would create different policy and persistence behavior.

### Practice guide

Implement:

```text
agentkit serve
agentkit chat --agent <id> --session demo "hello"
agentkit migrate up
agentkit providers verify
```

`serve` composes the application. `chat` sends an authenticated HTTP request and prints the final answer. `providers verify` sends a minimal provider request without exposing credentials.

## Task 6 — Verify the two adapters

### Practice guide

Test:

1. `/health` reports process and database readiness separately.
2. Unauthorized chat returns `401`.
3. Invalid JSON and oversized input return `400`/`413`.
4. A provider timeout returns a safe gateway error.
5. CLI and direct HTTP continue the same session history.
6. Client cancellation reaches the queued or active run.

This step is complete when HTTP and CLI produce the same durable agent behavior.
