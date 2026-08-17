# Step 6 — Agent Orchestration Loop

**Knowledge depth: 8/10**

This step combines prompt building, provider calls, tools, and persistence into one bounded agent run.

## Task 1 — Define per-run state

### Theory

Read [01 — Agent Loop](../01-agent-loop.md). GoClaw describes eight logical phases:

```text
context → history → prompt → think → act → observe → memory → summarize
```

Some phases can share code, but their responsibilities must remain visible. Run state belongs to one request; a reusable runner holds dependencies only.

### Practice guide

Create a private state type:

```go
type runState struct {
	request    RunRequest
	messages   []model.Message
	usage      model.Usage
	iterations int
	toolCalls  int
}

type Runner struct {
	Provider model.Provider
	Agents   AgentStore
	Sessions SessionStore
	Prompts  PromptBuilder
	Tools    ToolExecutor
	Config   Config
}
```

Do not reuse `runState` across calls. Do not store the current session or user on `Runner`.

## Task 2 — Prepare and persist the turn

### Practice guide

At the start of `Run`:

1. Validate agent, user, session key, and non-empty message.
2. Load agent settings.
3. Load scoped history and summary.
4. Select skills and available tool definitions.
5. Build canonical prompt messages.
6. Persist the user's message before the provider call.

If preparation fails, wrap the error with the phase name. Do not create a fictional assistant reply.

## Task 3 — Implement the bounded think loop

### Theory

The model either returns a final answer or proposes tool calls. The runtime executes allowed tools, appends observations, and asks the model again.

### Practice guide

Implement the central loop:

```go
for state.iterations < r.Config.MaxIterations {
	state.iterations++

	resp, err := r.Provider.Chat(ctx, model.ChatRequest{
		Model:     agentConfig.Model,
		Messages:  state.messages,
		Tools:     r.Tools.Definitions(ctx),
		MaxTokens: r.Config.MaxOutputTokens,
	})
	if err != nil {
		return resultFrom(state), fmt.Errorf("think iteration %d: %w", state.iterations, err)
	}

	state.usage = addUsage(state.usage, resp.Usage)
	assistant := model.Message{
		Role: "assistant", Content: resp.Content, ToolCalls: resp.ToolCalls,
	}
	state.messages = append(state.messages, assistant)

	if len(resp.ToolCalls) == 0 {
		// Persist the assistant message and return the final result.
	}
	// Persist the tool request, execute calls, and continue.
}
```

Apply context budgeting before every provider call; Step 7 implements the full pruning policy.

## Task 4 — Execute and observe tool calls

### Theory

The model proposes actions. The tool registry independently validates and authorizes them. Unknown, invalid, and denied calls become canonical tool error messages so the model may recover.

### Practice guide

For every tool-call batch:

1. Check the remaining total tool-call budget.
2. Persist the assistant message containing the original calls.
3. Execute each call through the registry.
4. Append and persist one tool message for each call ID.
5. Preserve provider order even if read-only I/O runs concurrently.
6. Continue to the next model iteration.

Never keep an assistant tool request without its matching tool results. Never execute the same call twice because persistence returned an ambiguous error.

## Task 5 — Define failure and cancellation behavior

### Theory

Read the run-visibility concepts in [10 — Tracing & Observability](../10-tracing-observability.md). A clear failure is better than a fluent but false assistant message.

### Practice guide

Implement these outcomes:

| Condition | Required behavior |
|---|---|
| Caller cancelled | Stop provider and tool work; perform safe final bookkeeping only. |
| Provider timeout | Return a retryable run error; do not persist an invented answer. |
| Unknown or denied tool | Append a tool error and let the model continue within limits. |
| Tool panic | Recover inside the registry and return a tool error. |
| Iteration/tool budget reached | Return a clear limit error with partial usage metadata. |
| Final persistence failed | Report failure; do not claim the turn is durable. |

Use `context.WithoutCancel` only for bounded, safe persistence. Do not continue model calls after cancellation.

## Task 6 — Test the complete loop

### Practice guide

Use scripted fake provider responses to test:

1. One text-only response.
2. One tool call followed by a final response.
3. Unknown tool followed by model recovery.
4. Repeated calls stopped by the tool-call limit.
5. Provider cancellation.
6. Persistence failure at each append point.

Assert message order, tool-call IDs, usage totals, and iteration counts. This step is complete when one runner can handle concurrent calls without sharing run state.
