# Step 5 — PostgreSQL Stores and Scope

**Knowledge depth: 7/10**

Implement scoped PostgreSQL adapters for sessions and agent settings.

## Step outcome

Session history and agent settings survive restart, preserve order, and reject cross-user reads.

## Task 1 — Implement PostgreSQL session storage

### Theory

- **Architectural role:** Transactions and locking belong to the PostgreSQL adapter. The agent sees only `Load` and `Append` behavior. The store protects its own invariant even before the scheduler from Step 19 exists.
- **Why:** Two writers can read the same `max(ordinal)` and create duplicates or incorrect order. Locking by session turns ordinal allocation and batch insertion into one database-level critical section.
- **GoClaw reference:** Inspect [`PGSessionStore`](../../goclaw/internal/store/pg/sessions.go), separated list queries in [`goclaw/internal/store/pg/sessions_list.go`](../../goclaw/internal/store/pg/sessions_list.go), and shared dialect/query helpers in [`goclaw/internal/store/base`](../../goclaw/internal/store/base). GoClaw uses a different data shape; study its transaction and scope responsibilities rather than copying its SQL exactly.

The transaction must preserve message order when multiple messages are appended together. Step 19 will serialize runs, but the store should still maintain its own invariant.

### Goal

Implement the `SessionStore` port while keeping a multi-message append atomic and ordered.

### Guide to implement

Implement `internal/store/postgres/session_store.go`:

1. Begin a transaction.
2. Upsert the scoped session row.
3. Lock the session row with `FOR UPDATE`.
4. Read the current maximum ordinal.
5. Insert messages with consecutive ordinals.
6. Update `sessions.updated_at` and commit.

`Load` must filter by `agent_id`, `user_id`, and `session_key`, order by ordinal, decode JSON into canonical messages, and return the stored summary.

Wrap errors with operation context such as `load session` or `append message`. Do not include full message content in errors.

## Task 2 — Add agent configuration storage

### Theory

- **Architectural role:** `AgentStore` is a separate aggregate repository. The runner reads a configuration snapshot at the start of a run, while transports and the UI update it through a validated use case.
- **Why:** Keeping settings only in a config file or browser can make different processes or users see different state. A patch type limits which fields callers may change instead of allowing arbitrary row updates.
- **GoClaw reference:** Read the contract in [`goclaw/internal/store/agent_store.go`](../../goclaw/internal/store/agent_store.go), the PostgreSQL adapter in [`goclaw/internal/store/pg/agents.go`](../../goclaw/internal/store/pg/agents.go), and gateway use cases in [`goclaw/internal/gateway/methods/agents_update.go`](../../goclaw/internal/gateway/methods/agents_update.go).

### Goal

Create the authoritative data source for an agent's model, prompt, skill, and tool settings.

### Guide to implement

Implement only the operations needed by later Steps:

```go
type AgentStore interface {
	Get(ctx context.Context, id string) (Agent, error)
	UpdateSettings(ctx context.Context, id string, patch AgentSettings) (Agent, error)
}
```

Seed one agent for local development. Its `system_prompt` and `model` will be consumed in Steps 8 and 10.

## Task 3 — Verify persistence and scope

### Theory

- **Architectural role:** The integration test exercises migrations, the adapter, and a real PostgreSQL instance. It catches constraints, transaction behavior, and query predicates that mocked SQL unit tests cannot verify.
- **Why:** Loading through a newly constructed store proves that state was not accidentally kept only in memory. Using the same `session_key` for another user proves that the key is not treated as authorization.
- **GoClaw reference:** Review scope tests such as [`goclaw/internal/store/pg/agents_list_tenant_scope_test.go`](../../goclaw/internal/store/pg/agents_list_tenant_scope_test.go), session-related tests under [`goclaw/internal/store/pg`](../../goclaw/internal/store/pg), and higher-level invariants under [`goclaw/tests/integration`](../../goclaw/tests/integration) when a PostgreSQL test environment is available.

### Goal

Verify two properties at the correct boundary: data survives a restart and cannot be read across scopes.

### Guide to implement

Write an integration test using a temporary PostgreSQL database or test schema:

1. Apply migrations.
2. Create one agent and two users.
3. Append a user and assistant message for user A.
4. Recreate the store instance and load the session.
5. Assert message order and content.
6. Load the same session key as user B and assert no messages are returned.

This step is complete when session history and agent settings survive a process/store restart and cross-user reads are impossible through the store API.
