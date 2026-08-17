# Step 16 — Chat UI, Run Timeline, and Project Review

**Knowledge depth: 8/10**

This step completes the real-time chat experience and verifies the full agent gateway vertical slice.

## Task 1 — Model chat and run states

### Theory

Read streaming flow in [01 — Agent Loop](../01-agent-loop.md), events in [19 — WebSocket RPC Methods](../19-websocket-rpc.md), and run visibility in [10 — Tracing & Observability](../10-tracing-observability.md).

Explicit states prevent stale events and visual races:

```text
composer: idle | submitting | streaming | awaiting_approval | cancelled
message: pending | complete | failed | interrupted
run: queued | running | completed | failed | cancelled
```

### Practice guide

Create a reducer or state machine keyed by `run_id` and `session_key`. Ignore or store late events that do not match the active session. A terminal event must close the correct pending assistant message exactly once.

## Task 2 — Build session navigation and composer

### Practice guide

Implement:

- Session list with create/select actions.
- Durable history loading on selection.
- Multiline composer with send and cancel actions.
- Disabled/clear status while a session turn is queued or running.
- Optimistic user message marked pending until request acceptance.

Generate a stable session key on the client or ask the server to create one. Never reuse a draft from one session after switching sessions.

## Task 3 — Render streaming messages safely

### Theory

Assistant text is untrusted output. Markdown rendering must sanitize HTML and generated links.

### Practice guide

On `chat.delta`, append text to the assistant message for the matching run. Batch frequent state updates if rendering becomes slow.

Render:

- Sanitized Markdown.
- Readable code blocks and tables.
- Safe links with appropriate target/rel behavior.
- Clear interrupted and failed states.

After `run.completed`, reload or reconcile with the durable server message so missed deltas cannot leave incorrect text.

## Task 4 — Show tools, memory, and approvals

### Practice guide

Display concise transcript activity:

```text
Reading workspace/report.md…
Searched 3 memory records
Write requires approval
```

Put full allowed arguments, redacted results, timing, and error details in an expandable panel. For approval-gated writes, show normalized action details and call the Step 14 approval endpoint. Never approve a changed argument set silently.

## Task 5 — Build the run timeline

### Theory

A timeline explains how an answer was produced without exposing hidden chain-of-thought. Show observable operations and outcomes only.

### Practice guide

Add a run-detail drawer containing:

- Queue wait and total duration.
- Provider model and token usage.
- Iteration and tool-call counts.
- Tool start/completion/error records.
- Memory retrieval count and source labels.
- Final status and safe error code.

Store or query these run facts from the gateway. Do not reconstruct the authoritative timeline only from WebSocket events.

## Task 6 — Implement reconnect recovery

### Practice guide

When the socket reconnects:

1. Mark active streams as reconnecting, not complete.
2. Reload the selected session messages.
3. Query any known active run IDs.
4. Reconcile pending UI items with durable state.
5. Resume listening for new events.

Test switching sessions during a reconnect and receiving a late event from the old session.

## Task 7 — Run the end-to-end acceptance check

### Practice guide

Complete this scenario:

1. Start PostgreSQL, apply migrations, and start the gateway/UI.
2. Edit the agent's system prompt and enable one local skill.
3. Create a session and send a message over WebSocket.
4. Observe text and one safe tool round-trip.
5. Complete another session so episodic memory is consolidated.
6. Start a new session and recall the scoped fact.
7. Restart the gateway and reload history/settings.
8. Open the run timeline and inspect usage, tools, and retrieval facts.
9. Attempt cross-scope access and a workspace traversal; both must fail.

## Task 8 — Review extension boundaries

### Theory

The remaining author documents explain platform features that should extend the core rather than rewrite it:

- [05 — Channels and Messaging](../05-channels-messaging.md)
- [11 — Agent Teams](../11-agent-teams.md)
- [14 — Skills Runtime](../14-skills-runtime.md)
- [15 — Core Skills System](../15-core-skills-system.md)
- [16 — Skill Publishing](../16-skill-publishing.md)
- [21 — Agent Evolution and Skill Management](../21-agent-evolution-and-skill-management.md)
- [22 — Heartbeat System](../22-heartbeat-system.md)

### Practice guide

For each future feature, name its extension point in the current project. Examples: a channel calls `ChatService`, a new provider implements `model.Provider`, cron submits a scheduler job, and publishing adds a managed skill store behind the existing skill loader.

The course is complete when the acceptance scenario passes and future features have clear extension points without partial implementations.
