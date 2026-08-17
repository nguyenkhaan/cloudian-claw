# Step 23 — Security Hardening, Approvals, and Audit

**Knowledge depth: 9/10**

Apply one trusted identity and defense-in-depth controls across every external and data boundary.

## Step outcome

Scoped reads, exact-action approvals, redacted audit events, and adversarial tests protect the complete single-tenant deployment.

## Task 1 — Write the threat model

### Theory

- **Architectural role:** The threat model connects data flow to defense in depth: edge authentication, scoped stores, tool authorization, output redaction, and audit. The prompt is only one soft control in the chain.
- **Why:** “Prompt injection protection” cannot replace path containment or SQL scope. Writing the flow makes it clear that retrieved data and tool output are also untrusted inputs.
- **GoClaw reference:** Start with [`goclaw/docs/09-security.md`](../../goclaw/docs/09-security.md), then inspect the input guard in [`goclaw/internal/agent/input_guard.go`](../../goclaw/internal/agent/input_guard.go), tool policy in [`goclaw/internal/tools/policy.go`](../../goclaw/internal/tools/policy.go), and the file-serving boundary in [`goclaw/internal/http/files.go`](../../goclaw/internal/http/files.go).

Read [09 — Security](../09-security.md), [20 — API Keys & Authentication](../20-api-keys-auth.md), and [23 — AI Agent Permission Matrix](../23-ai-agent-permission-matrix.md). Use [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) to understand future scope; do not build tenant administration or RBAC screens.

Prompt injection is not solved by a stronger system prompt. Security comes from scoped data access, least-privilege tools, validation, approvals, and audit logs.

### Goal

List trust boundaries, protected assets, attacker-controlled inputs, and enforcement points before adding controls.

### Guide to implement

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

- **Architectural role:** Credentials exist only at the authentication boundary, context carries an immutable principal, and resource selections from request bodies must be authorized against that principal.
- **Why:** Downstream token parsing or trust in payload identity creates inconsistent identities and bypass paths. Typed context helpers avoid string-key collisions and turn a missing principal into a clear error.
- **GoClaw reference:** Inspect identity context helpers in [`goclaw/internal/store/context.go`](../../goclaw/internal/store/context.go), actor tests in [`goclaw/internal/store/context_actor_id_test.go`](../../goclaw/internal/store/context_actor_id_test.go), HTTP authentication in [`goclaw/internal/http/auth.go`](../../goclaw/internal/http/auth.go), and WebSocket authentication in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go).

The request body can select a resource, but it cannot prove identity. Authentication middleware creates the trusted principal.

### Goal

Create one trusted principal used across HTTP, WebSocket, services, and stores.

### Guide to implement

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

### Theory

- **Architectural role:** The application authorizes intent, and store adapters always include trusted agent and user scope in SQL. Unauthorized and absent resources commonly return the same not-found result to prevent enumeration.
- **Why:** Loading broad data and filtering in Go exposes rows to memory or logs and makes it easy to forget a check at a new call site. Query-level scope reduces exposure at the lowest boundary.
- **GoClaw reference:** Study scope context in [`goclaw/internal/store/scope.go`](../../goclaw/internal/store/scope.go), shared tenant query helpers in [`goclaw/internal/store/base/tenant.go`](../../goclaw/internal/store/base/tenant.go), session queries in [`goclaw/internal/store/pg/sessions.go`](../../goclaw/internal/store/pg/sessions.go), and scope tests such as [`goclaw/internal/store/pg/agents_list_tenant_scope_test.go`](../../goclaw/internal/store/pg/agents_list_tenant_scope_test.go).

### Goal

Make ownership a required predicate at the point where data is read or written.

### Guide to implement

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

### Theory

- **Architectural role:** HTTP and WebSocket limit identity, input, and lifecycle; SQL limits queries and transactions; files limit namespaces; memory limits scope and trust; logging limits disclosure.
- **Why:** Each boundary has a different attack class. Rate limits do not stop traversal, parameterized SQL does not stop cross-user reads, and CORS does not replace authentication.
- **GoClaw reference:** Inspect the rate limiter in [`goclaw/internal/gateway/ratelimit.go`](../../goclaw/internal/gateway/ratelimit.go), origin and frame handling in [`goclaw/internal/gateway/server.go`](../../goclaw/internal/gateway/server.go), file-path tests in [`goclaw/internal/http/files_path_security_test.go`](../../goclaw/internal/http/files_path_security_test.go), and provider URL and SSRF validation in [`goclaw/internal/http/providers_url_validate_test.go`](../../goclaw/internal/http/providers_url_validate_test.go).

### Goal

Apply the right controls at each boundary instead of relying on one general “security” middleware.

### Guide to implement

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

- **Architectural role:** Normalized arguments produce the proposal identity hash, the server stores pending approval, and the registry checks and consumes it immediately before mutation. The UI requests approval but cannot grant permission by itself.
- **Why:** A simple `approved` boolean can be reused for a different path or content, or replayed. Binding actor, session, tool, and arguments means even a small change requires new consent.
- **GoClaw reference:** Find approval manager types with `rg -n 'type .*Approval|ApprovalManager' goclaw/internal/tools`, inspect the gateway method in [`goclaw/internal/gateway/methods/exec_approval.go`](../../goclaw/internal/gateway/methods/exec_approval.go), and follow registry and tool integration under [`goclaw/internal/tools`](../../goclaw/internal/tools).

Approval is a structured authorization fact, not a conversational “yes.” Bind it to the exact proposed action.

### Goal

Turn human approval into a one-time capability with scope, expiry, and binding to the exact action.

### Guide to implement

For `write_file`, model a pending approval with:

```text
actor + agent + session + tool name
+ hash(normalized arguments) + expiry
```

Normalize JSON deterministically before hashing. Consume approval once. If arguments change, require a new approval.

The first UI may use a confirm action that calls an approval endpoint. The tool registry must verify the server-side approval record before execution.

## Task 6 — Add audit events and redaction

### Theory

- **Architectural role:** Enforcement points emit structured events with correlation IDs, decisions, and reason codes. Redaction is a last-line guard before the sink, not a replacement for avoiding secret logging in the first place.
- **Why:** Audit needs to answer who requested what and how it was decided; it does not need complete prompts or tokens. Stable reason codes support alerts, metrics, and tests.
- **GoClaw reference:** Inspect the hook audit pattern in [`goclaw/internal/hooks/audit.go`](../../goclaw/internal/hooks/audit.go), gateway audit methods in [`goclaw/internal/gateway/methods/audit.go`](../../goclaw/internal/gateway/methods/audit.go), and scrubbing and redaction in [`goclaw/internal/tools/registry.go`](../../goclaw/internal/tools/registry.go).

### Goal

Record security decisions for investigation without turning the audit log into a store of secrets or content.

### Guide to implement

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

### Theory

- **Architectural role:** Every test identifies the control that blocks the input and also asserts non-effects: no row, file, or tool action occurs, and no secret reaches output or logs.
- **Why:** Checking only a status code can hide a side effect that happened before the error was returned. Cross-boundary tests catch authorization drift between HTTP, WebSocket, stores, and tools.
- **GoClaw reference:** Review authentication tests in [`goclaw/internal/http/auth_test.go`](../../goclaw/internal/http/auth_test.go), file security tests in [`goclaw/internal/http/files_path_security_test.go`](../../goclaw/internal/http/files_path_security_test.go), the API-key tenant guard in [`goclaw/internal/gateway/methods/api_keys_tenant_guard_test.go`](../../goclaw/internal/gateway/methods/api_keys_tenant_guard_test.go), and the invariant suite under [`goclaw/tests`](../../goclaw/tests).

### Goal

Verify defense in depth with adversarial tests that pass through public boundaries.

### Guide to implement

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
