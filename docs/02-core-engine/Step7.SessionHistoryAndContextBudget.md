# Step 7 — Session History and Context Budget

**Knowledge depth: 7/10**  

Read the history and compaction sections in [01 — Agent Loop](../01-agent-loop.md), then connect them to the persistence model in [06 — Store Layer and Data Model](../06-store-data-model.md). [07 — Bootstrap, Skills & Memory](../07-bootstrap-skills-memory.md) explains why a session summary is not long-term memory.

## Session is not memory

A **session** is the transcript for one active conversation. It optimizes continuity and auditability. It is not a good long-term knowledge store: it grows without limit, contains accidental details, and has only one conversation scope.

```text
Session history: exact recent interaction
Session summary: compressed older interaction
Long-term memory: retrieved facts from many sessions     ← Step 8
```

GoClaw has a focused `SessionStore` interface in `internal/store/session_store.go`, with PostgreSQL and SQLite implementations. Its pipeline loads history and summary before prompt construction, checks token pressure each iteration, checkpoints messages, and may summarize after a turn.

## Minimum schema

Store whole canonical messages as JSON initially; normalize later only when query needs justify it.

```sql
CREATE TABLE sessions (
  agent_id UUID NOT NULL,
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
  FOREIGN KEY (agent_id, user_id, session_key) REFERENCES sessions(agent_id, user_id, session_key)
);
CREATE INDEX session_messages_lookup
  ON session_messages (agent_id, user_id, session_key, ordinal);
```

Every query includes `agent_id` and `user_id`, even if the session key appears globally unique. GoClaw adds `tenant_id` for a multi-tenant deployment; Step 14 explains that extension without making it part of this project.

## Correct history construction

```go
func buildPrompt(system, summary string, history []model.Message, user string) []model.Message {
	msgs := []model.Message{{Role: "system", Content: system}}
	if summary != "" {
		msgs = append(msgs, model.Message{Role: "system", Content: "Conversation summary:\n" + summary})
	}
	msgs = append(msgs, history...)
	return append(msgs, model.Message{Role: "user", Content: user})
}
```

Rules:

1. Preserve tool-call pairs. Never keep an assistant request without its tool results, or a result without its request.
2. Append messages in provider-valid order.
3. Save the user input before model work if you want cancelled work to remain visible; otherwise use a transaction with an explicit status.
4. Do not use `[]Message` in an in-process map as your production source of truth.

## Budget and compaction policy

Use a real tokenizer for your model family in production. GoClaw centralizes this in `internal/tokencount` and tracks context window / output limits in pipeline configuration.

```go
func budget(window, systemTokens, toolTokens, outputLimit int) int {
	reserve := max(window/10, 512) // protect against estimation drift
	return max(0, window-systemTokens-toolTokens-outputLimit-reserve)
}
```

Prune in this order:

```text
1. remove transient UI events and obsolete scratch data
2. truncate oversized tool output, retaining an explicit truncation notice
3. preserve recent complete turns and all required tool-call pairs
4. summarize the oldest coherent block
5. save the summary, then remove the summarized raw history
```

The summary must be a durable checkpoint, not a string sent only to the next request. A crash after a transient compaction otherwise loses context forever.

Once this foundation is clear, the next two steps can add durable memory without confusing it with ordinary conversation history.
