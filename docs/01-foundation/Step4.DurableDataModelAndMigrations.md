# Step 4 — Durable Data Model and Migrations

**Knowledge depth: 6/10**

Model durable ownership and encode it in the first PostgreSQL migration.

## Step outcome

Agents, sessions, and ordered messages have explicit ownership, constraints, and indexes.

## Task 1 — Model durable ownership

### Theory

- **Architectural role:** An `agent` stores behavior settings, a `session` defines a conversation boundary by agent, user, and key, and a `message` is an ordered record within a session. This durable model is separate from temporary connection and run state.
- **Why:** Incorrect ownership leads to ambiguous queries and data leaks. A `session_key` identifies a conversation but does not prove that the caller may read it.
- **GoClaw reference:** Start with the domain contracts in [`goclaw/internal/store/agent_store.go`](../../goclaw/internal/store/agent_store.go) and [`goclaw/internal/store/session_store.go`](../../goclaw/internal/store/session_store.go), then read the dual-identity conventions in [`goclaw/docs/agent-identity-conventions.md`](../../goclaw/docs/agent-identity-conventions.md). This course simplifies multi-tenancy but keeps the ownership invariant.

Read [06 — Store Layer and Data Model](../06-store-data-model.md). Use [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) only to understand ownership; this project has one deployment scope.

The durable hierarchy is:

```text
project → agent → user → session → message
```

A session key is an identifier, not authorization. Every session query must also use agent and user scope.

### Goal

Decide which aggregates survive a restart and who owns them before designing tables.

### Guide to implement

Start with three aggregates:

- `agents`: model, editable system prompt, and enabled tool/skill settings.
- `sessions`: owner, conversation key, and durable summary.
- `session_messages`: ordered canonical messages.

Do not store connections, raw API keys, or large base64 payloads in these tables.

## Task 2 — Create the first migration

### Theory

- **Architectural role:** Migrations are the version history of the schema. Primary keys, foreign keys, and unique constraints protect relationships even when application code has a bug, while indexes support the store's actual query paths.
- **Why:** Checking uniqueness or ordering only in Go is not enough when multiple processes or transactions are active. Append-only migrations also make deployment and rollback behavior predictable.
- **GoClaw reference:** Review the migration sequence in [`goclaw/migrations`](../../goclaw/migrations), then find the session and agent tables with `rg 'CREATE TABLE.*sessions|CREATE TABLE.*agents' goclaw/migrations`. Compare them with consumer queries in [`goclaw/internal/store/pg/sessions.go`](../../goclaw/internal/store/pg/sessions.go) before choosing indexes.

### Goal

Encode the durable model and important invariants as database constraints and indexes.

### Guide to implement

Create a numbered migration in `migrations/`:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE agents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  model TEXT NOT NULL,
  system_prompt TEXT NOT NULL DEFAULT '',
  settings JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
  agent_id UUID NOT NULL REFERENCES agents(id),
  user_id TEXT NOT NULL,
  session_key TEXT NOT NULL,
  summary TEXT NOT NULL DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (agent_id, user_id, session_key)
);

CREATE TABLE session_messages (
  id BIGSERIAL PRIMARY KEY,
  agent_id UUID NOT NULL,
  user_id TEXT NOT NULL,
  session_key TEXT NOT NULL,
  ordinal BIGINT NOT NULL,
  message JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agent_id, user_id, session_key, ordinal),
  FOREIGN KEY (agent_id, user_id, session_key)
    REFERENCES sessions(agent_id, user_id, session_key)
);

CREATE INDEX session_messages_lookup
  ON session_messages (agent_id, user_id, session_key, ordinal);
```

Provide a matching down migration for local development. Never rewrite an applied production migration; add a new one.
