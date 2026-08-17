# Step 14 — Session History and Context Budget

**Knowledge depth: 8/10**

Build provider-valid history and fit it into a measured context budget.

## Step outcome

Long sessions preserve recent turns and complete tool pairs while staying within the model context window.

## Task 1 — Separate history, summary, and memory

### Theory

- **Architectural role:** History is the exact event log of a session, the summary is a lossy checkpoint of an older prefix, and long-term memory is a retrieval result that may come from another session. The prompt builder composes them without treating them as the same thing.
- **Why:** A vector store does not preserve tool protocol or order, a summary cannot replace the newest transcript, and retrieved memory is not trusted enough to become an instruction.
- **GoClaw reference:** Read the session history and summary contract in [`goclaw/internal/store/session_store.go`](../../goclaw/internal/store/session_store.go), history assembly in [`goclaw/internal/agent/loop_history.go`](../../goclaw/internal/agent/loop_history.go), and separate episodic injection in [`goclaw/internal/memory/auto_injector_impl.go`](../../goclaw/internal/memory/auto_injector_impl.go).

Read [01 — Agent Loop](../01-agent-loop.md), [06 — Store Layer and Data Model](../06-store-data-model.md), and the context sections of [07 — Bootstrap, Skills & Memory](../07-bootstrap-skills-memory.md).

These are different data layers:

```text
Recent history   exact messages for the current session
Session summary  compact checkpoint of older turns
Long-term memory retrieved facts across sessions (Steps 16–18)
```

Do not treat an embedding store as chat history. Provider protocols require exact message and tool-call ordering.

### Goal

Understand three context types with different provenance, precision, and lifecycles.

### Guide to implement

Review the Step 4 schema and Step 5 store. Confirm that `Load` returns summary separately from ordered messages. Add a query that can load a bounded recent suffix without changing message ordinals.

## Task 2 — Build provider-valid history

### Theory

- **Architectural role:** The history builder is the normalization point between durable messages and provider requests. It inserts system content, summaries, and the current turn while preserving complete tool groups.
- **Why:** Building messages at multiple call sites can duplicate the current user message or create orphaned tool results. A WebSocket delta is a presentation event, not a durable conversation fact.
- **GoClaw reference:** Follow `buildMessages` in [`goclaw/internal/agent/loop_history.go`](../../goclaw/internal/agent/loop_history.go), message validation and sanitization in [`goclaw/internal/agent/sanitize.go`](../../goclaw/internal/agent/sanitize.go), and the canonical buffer in [`goclaw/internal/pipeline/message_buffer.go`](../../goclaw/internal/pipeline/message_buffer.go).

### Goal

Build one canonical transcript that satisfies provider protocol before every model call.

### Guide to implement

Implement history construction in one place:

```go
func buildMessages(system, summary string, history []model.Message, user string) []model.Message {
	msgs := []model.Message{{Role: "system", Content: system}}
	if summary != "" {
		msgs = append(msgs, model.Message{
			Role: "system", Content: "Conversation summary:\n" + summary,
		})
	}
	msgs = append(msgs, history...)
	return append(msgs, model.Message{Role: "user", Content: user})
}
```

Enforce these invariants:

1. Roles appear in a provider-valid order.
2. An assistant tool request and all matching tool results stay together.
3. The current user message appears once.
4. Transient WebSocket events never enter provider history.

## Task 3 — Count and allocate tokens

### Theory

- **Architectural role:** `TokenCounter` is a tokenizer port. The allocator reserves fixed or estimated space for system content, tool schemas, output, and a safety margin before giving the remainder to history.
- **Why:** Tool definitions and model output also consume context. A counter abstraction can use a real tokenizer when available and a conservative fallback for unsupported models.
- **GoClaw reference:** Inspect the interfaces and implementations in [`goclaw/internal/tokencount/token_counter.go`](../../goclaw/internal/tokencount/token_counter.go) and [`goclaw/internal/tokencount/tiktoken_counter.go`](../../goclaw/internal/tokencount/tiktoken_counter.go), budget logic in [`goclaw/internal/tokencount/budget_counter.go`](../../goclaw/internal/tokencount/budget_counter.go), and request budgeting in [`goclaw/internal/agent/loop_request_budget.go`](../../goclaw/internal/agent/loop_request_budget.go).

Context space must reserve room for the system prompt, tool schemas, model output, and estimation error. Use a tokenizer compatible with the chosen model when available.

### Goal

Treat the context window as a budget allocated in advance instead of trimming only after the provider rejects a request.

### Guide to implement

Centralize token counting behind an interface:

```go
type TokenCounter interface {
	CountMessages([]model.Message) int
	CountText(string) int
}
```

Calculate the history budget:

```go
func historyBudget(window, systemTokens, toolTokens, outputLimit int) int {
	reserve := max(window/10, 512)
	return max(0, window-systemTokens-toolTokens-outputLimit-reserve)
}
```

Fail early if stable system instructions and tool definitions already exceed the window.

## Task 4 — Prune without breaking tool pairs

### Theory

- **Architectural role:** Pruning works on semantic groups or turns rather than individual message indexes. Oversized observations are replaced with explicit markers before older turns are removed.
- **Why:** An assistant tool call and its result have a referential relationship. Removing half of the pair gives the model an unknown ID or makes an action appear unfinished. A marker tells the model that data was removed because of the budget.
- **GoClaw reference:** Study policies in [`goclaw/internal/agent/pruning.go`](../../goclaw/internal/agent/pruning.go), integration behavior in [`goclaw/internal/agent/pruning_integration_test.go`](../../goclaw/internal/agent/pruning_integration_test.go), and separate tool-result truncation in [`goclaw/internal/agent/tool_result_truncation.go`](../../goclaw/internal/agent/tool_result_truncation.go).

### Goal

Reduce context while preserving meaningful, provider-valid conversation structure.

### Guide to implement

Apply this order:

1. Remove transient scratch content that is not durable conversation.
2. Truncate oversized tool output and add an explicit truncation notice.
3. Preserve the newest complete turns.
4. Select the oldest complete turns for summarization.
5. Never split an assistant tool request from its results.

Represent complete turns or message groups in code instead of deleting individual messages by index.
