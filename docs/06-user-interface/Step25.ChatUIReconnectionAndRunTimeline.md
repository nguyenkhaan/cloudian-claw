# Step 25 — Chat UI, Reconnection, and Run Timeline

**Knowledge depth: 8/10**

Complete real-time chat presentation, tool approval UX, recovery, and observable run details.

## Step outcome

Users can chat, inspect safe tool and memory activity, reconnect without corrupting state, and view a durable run timeline.

## Task 1 — Model chat and run states

### Theory

- **Architectural role:** A reducer or state machine is keyed by `run_id + session_key`, transitions accept only valid events, and terminal states close the correct pending message idempotently.
- **Why:** Separate booleans such as `loading`, `streaming`, and `error` can create impossible combinations. Correlation keys stop late events from an old session from changing the new view.
- **GoClaw reference:** Inspect the chat-message projection store in [`goclaw/ui/web/src/stores/use-chat-messages-store.ts`](../../goclaw/ui/web/src/stores/use-chat-messages-store.ts), send flow in [`goclaw/ui/web/src/pages/chat/hooks/use-chat-send.ts`](../../goclaw/ui/web/src/pages/chat/hooks/use-chat-send.ts), and event hook in [`goclaw/ui/web/src/hooks/use-ws-event.ts`](../../goclaw/ui/web/src/hooks/use-ws-event.ts).

Read streaming flow in [01 — Agent Loop](../01-agent-loop.md), events in [19 — WebSocket RPC Methods](../19-websocket-rpc.md), and run visibility in [10 — Tracing & Observability](../10-tracing-observability.md).

Explicit states prevent stale events and visual races:

```text
composer: idle | submitting | streaming | awaiting_approval | cancelled
message: pending | complete | failed | interrupted
run: queued | running | completed | failed | cancelled
```

### Goal

Represent chat and run lifecycles explicitly so asynchronous events do not create UI races.

### Guide to implement

Create a reducer or state machine keyed by `run_id` and `session_key`. Ignore or store late events that do not match the active session. A terminal event must close the correct pending assistant message exactly once.

## Task 2 — Build session navigation and composer

### Theory

- **Architectural role:** The URL or selection identifies a session, a query loads authoritative history, the composer keeps a local draft, and an optimistic message remains a stateful projection until server acknowledgement.
- **Why:** Draft and pending state are not durable facts. Resetting them by session prevents sending the wrong context, while disabled and cancel states reflect scheduler and run state rather than only the HTTP request.
- **GoClaw reference:** Inspect [`goclaw/ui/web/src/pages/chat/chat-sidebar.tsx`](../../goclaw/ui/web/src/pages/chat/chat-sidebar.tsx), the session hook in [`goclaw/ui/web/src/pages/chat/hooks/use-chat-sessions.ts`](../../goclaw/ui/web/src/pages/chat/hooks/use-chat-sessions.ts), input in [`goclaw/ui/web/src/components/chat/chat-input.tsx`](../../goclaw/ui/web/src/components/chat/chat-input.tsx), and the session-key helper in [`goclaw/ui/web/src/lib/session-key.ts`](../../goclaw/ui/web/src/lib/session-key.ts).

### Goal

Connect session identity, durable history, and input lifecycle into the first complete user flow.

### Guide to implement

Implement:

- Session list with create/select actions.
- Durable history loading on selection.
- Multiline composer with send and cancel actions.
- Disabled/clear status while a session turn is queued or running.
- Optimistic user message marked pending until request acceptance.

Generate a stable session key on the client or ask the server to create one. Never reuse a draft from one session after switching sessions.

## Task 3 — Render streaming messages safely

### Theory

- **Architectural role:** The event reducer appends raw text by correlation key, the presentation renderer sanitizes and formats it, and completion triggers reconciliation with the server. Batching is a performance concern and does not change semantics.
- **Why:** Model output contains untrusted HTML and links. Deltas may be missed during disconnects; only a durable reload confirms final content and prevents the UI from keeping an incorrectly assembled message.
- **GoClaw reference:** Inspect the streaming renderer in [`goclaw/ui/web/src/components/chat/streaming-text.tsx`](../../goclaw/ui/web/src/components/chat/streaming-text.tsx), the content boundary in [`goclaw/ui/web/src/components/chat/message-content.tsx`](../../goclaw/ui/web/src/components/chat/message-content.tsx), the rich parser in [`goclaw/ui/web/src/components/chat/rich-content-parser.ts`](../../goclaw/ui/web/src/components/chat/rich-content-parser.ts), and the message hook in [`goclaw/ui/web/src/pages/chat/hooks/use-chat-messages.ts`](../../goclaw/ui/web/src/pages/chat/hooks/use-chat-messages.ts).

Assistant text is untrusted output. Markdown rendering must sanitize HTML and generated links.

### Goal

Render deltas quickly and safely, then converge on the final durable message.

### Guide to implement

On `chat.delta`, append text to the assistant message for the matching run. Batch frequent state updates if rendering becomes slow.

Render:

- Sanitized Markdown.
- Readable code blocks and tables.
- Safe links with appropriate target/rel behavior.
- Clear interrupted and failed states.

After `run.completed`, reload or reconcile with the durable server message so missed deltas cannot leave incorrect text.

## Task 4 — Show tools, memory, and approvals

### Theory

- **Architectural role:** The transcript shows a safe summary, expandable details use redacted observable data, and approval mutations send the exact normalized proposal ID or hash to the server.
- **Why:** The UI must not infer approval from the word “yes” or from changed arguments. Progressive disclosure keeps chat readable while still allowing users to audit actions.
- **GoClaw reference:** Inspect [`goclaw/ui/web/src/components/chat/tool-call-card.tsx`](../../goclaw/ui/web/src/components/chat/tool-call-card.tsx), find approval flows with `rg -n 'approval' goclaw/ui/web/src`, review memory types in [`goclaw/ui/web/src/types/memory.ts`](../../goclaw/ui/web/src/types/memory.ts), and inspect the backend approval method in [`goclaw/internal/gateway/methods/exec_approval.go`](../../goclaw/internal/gateway/methods/exec_approval.go).

### Goal

Make side effects and retrieval observable and consent-based without turning the transcript into a debug dump.

### Guide to implement

Display concise transcript activity:

```text
Reading workspace/report.md…
Searched 3 memory records
Write requires approval
```

Put full allowed arguments, redacted results, timing, and error details in an expandable panel. For approval-gated writes, show normalized action details and call the Step 23 approval endpoint. Never approve a changed argument set silently.

## Task 5 — Build the run timeline

### Theory

- **Architectural role:** The recorder and store are the authoritative append and query path, WebSocket only updates a live projection, and the drawer reads a durable timeline of stages, tools, usage, and status.
- **Why:** Reconstructing the timeline only from events loses data during disconnects and can differ between clients. Operational metadata is enough for debugging without reasoning tokens or sensitive content.
- **GoClaw reference:** Inspect the recorder in [`goclaw/internal/agent/run_timeline_recorder.go`](../../goclaw/internal/agent/run_timeline_recorder.go), store contract in [`goclaw/internal/store/run_timeline_store.go`](../../goclaw/internal/store/run_timeline_store.go), PostgreSQL adapter in [`goclaw/internal/store/pg/run_timeline.go`](../../goclaw/internal/store/pg/run_timeline.go), and gateway query method in [`goclaw/internal/gateway/methods/run_timeline.go`](../../goclaw/internal/gateway/methods/run_timeline.go).

A timeline explains how an answer was produced without exposing hidden chain-of-thought. Show observable operations and outcomes only.

### Goal

Explain a run with observable facts without exposing hidden chain of thought.

### Guide to implement

Add a run-detail drawer containing:

- Queue wait and total duration.
- Provider model and token usage.
- Iteration and tool-call counts.
- Tool start/completion/error records.
- Memory retrieval count and source labels.
- Final status and safe error code.

Store or query these run facts from the gateway. Do not reconstruct the authoritative timeline only from WebSocket events.

## Task 6 — Implement reconnect recovery

### Theory

- **Architectural role:** Reconnect state does not mark success or failure by itself. The client reloads the session and known runs, merges by IDs, and then receives new events. Correlation filters events from old sessions.
- **Why:** A disconnect does not mean the server run stopped. Closing a pending message early creates a false terminal state, while appending an undeduplicated snapshot creates duplicates.
- **GoClaw reference:** Follow reconnection in [`goclaw/ui/web/src/api/ws-client.ts`](../../goclaw/ui/web/src/api/ws-client.ts), the WebSocket provider in [`goclaw/ui/web/src/components/providers/ws-provider.tsx`](../../goclaw/ui/web/src/components/providers/ws-provider.tsx), the chat-message store in [`goclaw/ui/web/src/stores/use-chat-messages-store.ts`](../../goclaw/ui/web/src/stores/use-chat-messages-store.ts), and active-run UI in [`goclaw/ui/web/src/components/chat/active-run-zone.tsx`](../../goclaw/ui/web/src/components/chat/active-run-zone.tsx).

### Goal

Reconcile optimistic and streaming UI with durable messages and run state after a connection gap.

### Guide to implement

When the socket reconnects:

1. Mark active streams as reconnecting, not complete.
2. Reload the selected session messages.
3. Query any known active run IDs.
4. Reconcile pending UI items with durable state.
5. Resume listening for new events.

Test switching sessions during a reconnect and receiving a late event from the old session.
