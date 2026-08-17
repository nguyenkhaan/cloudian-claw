# Step 14 — Single-Tenant Security and Access Control

**Knowledge depth: 8/10**  

Before this stage, study [09 — Security](../09-security.md) and [20 — API Keys & Authentication](../20-api-keys-auth.md). Read [23 — AI Agent Permission Matrix](../23-ai-agent-permission-matrix.md) and [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) to understand how GoClaw scales authority and scope. This project remains a single tenant with one configured access token; it does not build tenant administration or RBAC screens.

## Threat model first

An agent combines untrusted input with powerful actions. Protect the following flows:

```text
user input → prompt → model proposal → tool call → filesystem/network/DB
retrieved memory/document → prompt
browser/API client → HTTP or WebSocket → user-scoped stores
configured credential → provider request
```

Prompt injection is not solved by a system prompt. It is reduced through least privilege, tool policy, scoped retrieval, output validation, and auditable approvals.

## A single project still needs scope

Make ownership impossible to forget in a store query. GoClaw injects tenant, agent, user, and locale into `context.Context`; its tenant-aware stores then add predicates. For this course, propagate only the user and agent identity, and keep it in every session, message, memory, and trace query.

```go
type principalKey struct{}
type Principal struct{ UserID, AgentID string }
func WithPrincipal(ctx context.Context, p Principal) context.Context { return context.WithValue(ctx, principalKey{}, p) }
func PrincipalFrom(ctx context.Context) (Principal, error) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	if !ok || p.UserID == "" || p.AgentID == "" { return Principal{}, errors.New("principal required") }
	return p, nil
}

func (s *SessionStore) Load(ctx context.Context, key string) ([]model.Message, error) {
	p, err := PrincipalFrom(ctx); if err != nil { return nil, err }
	rows, err := s.db.QueryContext(ctx, `SELECT message FROM session_messages
WHERE user_id = $1 AND agent_id = $2 AND session_key = $3 ORDER BY ordinal`, p.UserID, p.AgentID, key)
	// scan rows…
	return messages, err
}
```

Never accept `user_id` or `agent_id` from a normal request body as authority. Resolve both from the authenticated connection or selected server-side agent. In GoClaw's multi-tenant architecture, this same rule expands to tenant scope and role checks.

## Boundary controls

| Boundary | Required control |
|---|---|
| HTTP | authentication, request size limit, rate limit, CORS, typed validation |
| WebSocket | authenticated handshake, origin check, frame/read limit, outbound backpressure |
| SQL | parameterized statements and user/agent scope predicates |
| Files | canonical path, workspace-root containment, deny sensitive paths |
| Network tools | allow-list or SSRF protections; block private/link-local metadata targets |
| Shell | not included in this project |
| Secrets | environment variables and redaction from logs/results |
| Memory | scope before search; retrieved text is reference data, never authority |

GoClaw has concrete counterparts: AES-256-GCM credential handling (`internal/crypto`), path controls in `internal/tools`, SSRF-aware web tooling, approval machinery, input guard, gateway rate limiting, and `internal/permissions` RBAC.

## Approval is a product flow

For a workspace write, either request simple confirmation in the UI or record a structured pending approval. Bind the approval to:

```text
actor + agent + session + tool name + normalized arguments + expiry
```

Hash the normalized arguments. If the model changes `rm report.txt` into `rm *`, the old approval must not apply.

## Audit events

Record decisions, not secrets:

```json
{"kind":"tool.authorization","user":"…","agent":"…","tool":"write_file","decision":"denied","reason":"approval_required","request_id":"…"}
```

Useful audit fields: request/run ID, actor, agent, tool capability, decision, timing, resource identifier, and a redacted argument summary.

The security model should be visible in the design of every earlier step, not added as a final wrapper around the finished agent.
