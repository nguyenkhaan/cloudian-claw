# Step 7 — Session History and Context Budget

**Knowledge depth: 8/10**

This step keeps long conversations valid, ordered, and within the model context window.

## Task 1 — Separate history, summary, and memory

### Theory

Read [01 — Agent Loop](../01-agent-loop.md), [06 — Store Layer and Data Model](../06-store-data-model.md), and the context sections of [07 — Bootstrap, Skills & Memory](../07-bootstrap-skills-memory.md).

These are different data layers:

```text
Recent history   exact messages for the current session
Session summary  compact checkpoint of older turns
Long-term memory retrieved facts across sessions (Steps 8–9)
```

Do not treat an embedding store as chat history. Provider protocols require exact message and tool-call ordering.

### Practice guide

Review the Step 3 schema and store. Confirm that `Load` returns summary separately from ordered messages. Add a query that can load a bounded recent suffix without changing message ordinals.

## Task 2 — Build provider-valid history

### Practice guide

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

Context space must reserve room for the system prompt, tool schemas, model output, and estimation error. Use a tokenizer compatible with the chosen model when available.

### Practice guide

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

### Practice guide

Apply this order:

1. Remove transient scratch content that is not durable conversation.
2. Truncate oversized tool output and add an explicit truncation notice.
3. Preserve the newest complete turns.
4. Select the oldest complete turns for summarization.
5. Never split an assistant tool request from its results.

Represent complete turns or message groups in code instead of deleting individual messages by index.

## Task 5 — Save a durable summary checkpoint

### Theory

A summary sent only to the next model call is unsafe. A process crash would lose the compacted context.

### Practice guide

When compaction is required:

1. Select the oldest coherent message range.
2. Ask the provider for a factual, instruction-free summary.
3. Store the new summary.
4. Mark or delete only the summarized range in the same transaction or a recoverable sequence.
5. Rebuild the prompt and confirm it fits the budget.

Include durable facts, decisions, unresolved questions, and important references. Exclude hidden reasoning and any claim not present in the source messages.

## Task 6 — Verify long-session behavior

### Practice guide

Create a fixture with many turns, at least two tool-call pairs, and one oversized tool result. Assert that:

- The final prompt is within budget.
- Recent turns remain exact.
- Tool pairs remain complete and ordered.
- Truncation is visible.
- The summary survives a new store instance.
- Repeated compaction does not summarize the same range twice.

This step is complete when a long session can continue after restart without invalid provider history.
