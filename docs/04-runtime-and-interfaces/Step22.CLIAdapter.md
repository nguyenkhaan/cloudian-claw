# Step 22 — CLI Adapter

**Knowledge depth: 6/10**

Expose operator and chat workflows through a thin gateway client.

## Step outcome

The CLI and direct HTTP path share authentication, persistence, policy, cancellation, and session behavior.

## Task 1 — Implement CLI commands

### Theory

- **Architectural role:** `serve` is a composition command, while `chat` and `verify` are remote adapters or clients. The CLI parses and formats; the server owns policy and durable behavior.
- **Why:** A second runner embedded in the CLI can bypass scheduling and authentication or use different configuration. An HTTP client keeps one source of truth and also tests the deployed gateway path.
- **GoClaw reference:** Inspect command definitions in [`goclaw/cmd/agent.go`](../../goclaw/cmd/agent.go), the chat client in [`goclaw/cmd/agent_chat_client.go`](../../goclaw/cmd/agent_chat_client.go), gateway startup in [`goclaw/cmd/gateway.go`](../../goclaw/cmd/gateway.go), and provider verification in [`goclaw/cmd/providers_cmd.go`](../../goclaw/cmd/providers_cmd.go).

The CLI is an HTTP client of the running gateway. Building a second in-process runner would create different policy and persistence behavior.

### Goal

Provide an operator and user CLI while correctly reusing the gateway contract.

### Guide to implement

Implement:

```text
agentkit serve
agentkit chat --agent <id> --session demo "hello"
agentkit migrate up
agentkit providers verify
```

`serve` composes the application. `chat` sends an authenticated HTTP request and prints the final answer. `providers verify` sends a minimal provider request without exposing credentials.

## Task 2 — Verify the two adapters

### Theory

- **Architectural role:** Adapter tests verify mapping and cancellation. An end-to-end test uses the same session key through both entry points to confirm one durable core.
- **Why:** Separate unit tests for a handler and client can pass even when they use different endpoints or schemas. A cross-adapter scenario catches contract drift.
- **GoClaw reference:** Review the HTTP contract in [`goclaw/internal/http/chat_completions.go`](../../goclaw/internal/http/chat_completions.go), the CLI client in [`goclaw/cmd/agent_chat_client.go`](../../goclaw/cmd/agent_chat_client.go), and gateway client tests in [`goclaw/internal/gateway/chat_runner_test.go`](../../goclaw/internal/gateway/chat_runner_test.go).

### Goal

Prove that HTTP and CLI differ only in presentation while sharing authentication, service, and session behavior.

### Guide to implement

Test:

1. `/health` reports process and database readiness separately.
2. Unauthorized chat returns `401`.
3. Invalid JSON and oversized input return `400`/`413`.
4. A provider timeout returns a safe gateway error.
5. CLI and direct HTTP continue the same session history.
6. Client cancellation reaches the queued or active run.

This step is complete when HTTP and CLI produce the same durable agent behavior.
