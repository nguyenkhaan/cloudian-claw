# Step 3 — Data Model and Store Layer

**Knowledge depth: 7/10**

This step makes agents, sessions, and messages survive process restarts.

## Task 1 — Model durable ownership

### Theory

Read [06 — Store Layer and Data Model](../06-store-data-model.md). Use [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) only to understand ownership; this project has one deployment scope.

The durable hierarchy is:

```text
project → agent → user → session → message
```

A session key is an identifier, not authorization. Every session query must also use agent and user scope.

### Practice guide

Start with three aggregates:

- `agents`: model, editable system prompt, and enabled tool/skill settings.
- `sessions`: owner, conversation key, and durable summary.
- `session_messages`: ordered canonical messages.

Do not store connections, raw API keys, or large base64 payloads in these tables.

## Task 2 — Create the first migration

### Practice guide

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

## Task 3 — Implement PostgreSQL session storage

### Theory

The transaction must preserve message order when multiple messages are appended together. Step 11 will serialize runs, but the store should still maintain its own invariant.

### Practice guide

Implement `internal/store/postgres/session_store.go`:

1. Begin a transaction.
2. Upsert the scoped session row.
3. Lock the session row with `FOR UPDATE`.
4. Read the current maximum ordinal.
5. Insert messages with consecutive ordinals.
6. Update `sessions.updated_at` and commit.

`Load` must filter by `agent_id`, `user_id`, and `session_key`, order by ordinal, decode JSON into canonical messages, and return the stored summary.

Wrap errors with operation context such as `load session` or `append message`. Do not include full message content in errors.

## Task 4 — Add agent configuration storage

### Practice guide

Implement only the operations needed by later Steps:

```go
type AgentStore interface {
	Get(ctx context.Context, id string) (Agent, error)
	UpdateSettings(ctx context.Context, id string, patch AgentSettings) (Agent, error)
}
```

Seed one agent for local development. Its `system_prompt` and `model` will be consumed in Step 5.

## Task 5 — Verify persistence and scope

### Practice guide

Write an integration test using a temporary PostgreSQL database or test schema:

1. Apply migrations.
2. Create one agent and two users.
3. Append a user and assistant message for user A.
4. Recreate the store instance and load the session.
5. Assert message order and content.
6. Load the same session key as user B and assert no messages are returned.

This step is complete when session history and agent settings survive a process/store restart and cross-user reads are impossible through the store API.
