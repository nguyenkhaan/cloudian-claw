# Step 16 — Tenant Isolation and Security Controls

**Knowledge depth: 8/10**  

Before this stage, study [09 — Security](../09-security.md) and [20 — API Keys & Authentication](../20-api-keys-auth.md). When you introduce multiple users or tenants, use [23 — AI Agent Permission Matrix](../23-ai-agent-permission-matrix.md) together with [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) to understand how GoClaw keeps authority and data scope aligned.

## Threat model first

An agent combines untrusted input with powerful actions. Protect the following flows:

```text
user/channel input → prompt → model proposal → tool call → filesystem/network/DB
retrieved memory/document → prompt
browser/API client → HTTP or WebSocket → tenant-scoped stores
encrypted credential → provider/MCP request
```

Prompt injection is not solved by a system prompt. It is reduced through least privilege, tool policy, scoped retrieval, output validation, and auditable approvals.

## Tenant context is an invariant

Make the scope impossible to forget in a store query. GoClaw injects tenant, agent, user, and locale into `context.Context`; its tenant-aware stores then add predicates. The same rule must apply to reads **and** writes.

```go
type tenantKey struct{}
func WithTenant(ctx context.Context, id uuid.UUID) context.Context { return context.WithValue(ctx, tenantKey{}, id) }
func Tenant(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(tenantKey{}).(uuid.UUID)
	if !ok || id == uuid.Nil { return uuid.Nil, errors.New("tenant scope required") }
	return id, nil
}

func (s *SessionStore) Load(ctx context.Context, key string) ([]model.Message, error) {
	tenant, err := Tenant(ctx); if err != nil { return nil, err }
	rows, err := s.db.QueryContext(ctx, `SELECT message FROM session_messages
WHERE tenant_id = $1 AND session_key = $2 ORDER BY ordinal`, tenant, key)
	// scan rows…
	return messages, err
}
```

Never accept `tenant_id` from a normal request body as authority. Resolve it from the authenticated connection/session. Cross-tenant administration requires a deliberate master-scope path and audit trail.

## Boundary controls

| Boundary | Required control |
|---|---|
| HTTP | authentication, request size limit, rate limit, CORS, typed validation |
| WebSocket | authenticated handshake, origin check, frame/read limit, outbound backpressure |
| SQL | parameterized statements, tenant predicate, indexes that include scope |
| Files | canonical path, workspace-root containment, deny sensitive paths |
| Network tools | allow-list or SSRF protections; block private/link-local metadata targets |
| Shell | disabled by default; sandbox, command deny rules, timeout, approval |
| Secrets | environment/key manager, encryption at rest, redaction from logs/results |
| Memory | scope before search; retrieved text is reference data, never authority |

GoClaw has concrete counterparts: AES-256-GCM credential handling (`internal/crypto`), path controls in `internal/tools`, SSRF-aware web tooling, approval machinery, input guard, gateway rate limiting, and `internal/permissions` RBAC.

## Approval is a product flow

For destructive or external actions, return a structured pending approval instead of executing. Bind the approval to:

```text
tenant + actor + agent + session + tool name + normalized arguments + expiry
```

Hash the normalized arguments. If the model changes `rm report.txt` into `rm *`, the old approval must not apply.

## Audit events

Record decisions, not secrets:

```json
{"kind":"tool.authorization","tenant":"…","agent":"…","tool":"write_file","decision":"denied","reason":"approval_required","request_id":"…"}
```

Useful audit fields: request/run ID, actor, tenant, agent, tool capability, decision, timing, resource identifier, and a redacted argument summary.

The security model should be visible in the design of every earlier step, not added as a final wrapper around the finished agent.
