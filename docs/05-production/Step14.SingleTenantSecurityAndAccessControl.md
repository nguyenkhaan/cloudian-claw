# Step 14 — Single-Tenant Security and Access Control

**Knowledge depth: 8/10**

This step applies one trusted identity and consistent security controls across HTTP, WebSocket, stores, memory, and tools.

## Task 1 — Write the threat model

### Theory

Read [09 — Security](../09-security.md), [20 — API Keys & Authentication](../20-api-keys-auth.md), and [23 — AI Agent Permission Matrix](../23-ai-agent-permission-matrix.md). Use [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) to understand future scope; do not build tenant administration or RBAC screens.

Prompt injection is not solved by a stronger system prompt. Security comes from scoped data access, least-privilege tools, validation, approvals, and audit logs.

### Practice guide

Create `docs/threat-model.md` for these paths:

```text
user input → prompt → model proposal → tool → file/memory
retrieved memory/document → prompt
browser/CLI → HTTP/WS → scoped stores
configured secret → provider request
```

For each path, list the protected asset, attacker-controlled input, enforcement point, and expected audit event.

## Task 2 — Resolve identity once at the edge

### Theory

The request body can select a resource, but it cannot prove identity. Authentication middleware creates the trusted principal.

### Practice guide

Create context helpers:

```go
type Principal struct {
	UserID         string
	AllowedAgentID string
}

func WithPrincipal(context.Context, Principal) context.Context
func PrincipalFrom(context.Context) (Principal, error)
```

Use the same authentication code for HTTP and WebSocket upgrade requests. Compare tokens in constant time and rotate them through configuration, not source code.

## Task 3 — Enforce scope in every store

### Practice guide

Audit session, message, memory, agent, and trace queries. Each query must derive user/agent scope from the principal or a server-side authorization decision.

Example:

```sql
SELECT message
FROM session_messages
WHERE user_id = $1
  AND agent_id = $2
  AND session_key = $3
ORDER BY ordinal;
```

Never implement “load by session key, then check owner in Go.” Scope in SQL so unauthorized rows are not loaded.

## Task 4 — Harden transport boundaries

### Practice guide

Apply:

| Boundary | Controls |
|---|---|
| HTTP | Authentication, body limit, rate limit, CORS, typed validation, timeouts. |
| WebSocket | Authenticated upgrade, origin check, frame limit, idle deadline, bounded write queue. |
| SQL | Parameters, scoped predicates, transaction deadlines. |
| Files | Canonical path and workspace containment. |
| Memory | Scope before ranking; retrieved content marked as data. |
| Secrets | Environment/config source and redaction. |

The project does not include shell or network tools. If added later, they require sandboxing and SSRF controls.

## Task 5 — Implement write approval

### Theory

Approval is a structured authorization fact, not a conversational “yes.” Bind it to the exact proposed action.

### Practice guide

For `write_file`, model a pending approval with:

```text
actor + agent + session + tool name
+ hash(normalized arguments) + expiry
```

Normalize JSON deterministically before hashing. Consume approval once. If arguments change, require a new approval.

The first UI may use a confirm action that calls an approval endpoint. The tool registry must verify the server-side approval record before execution.

## Task 6 — Add audit events and redaction

### Practice guide

Record decisions, not secrets:

```json
{
  "kind":"tool.authorization",
  "request_id":"...",
  "user_id":"...",
  "agent_id":"...",
  "tool":"write_file",
  "decision":"denied",
  "reason":"approval_required"
}
```

Audit authentication failures, tool allow/deny, approval creation/use, and cross-scope not-found results. Redact known tokens, authorization headers, and provider keys before log output.

## Task 7 — Run security checks

### Practice guide

Test:

1. JSON `user_id` cannot change the principal.
2. A user cannot read another user's session or memory.
3. Invalid WebSocket origin/token is rejected.
4. Traversal and symlink escape cannot leave the workspace.
5. A write without exact approval is denied.
6. Changed arguments invalidate approval.
7. Secrets do not appear in logs or errors.
8. Rate and input limits return stable errors.

This step is complete when all external paths use one trusted identity and every powerful action has an independent enforcement point.
