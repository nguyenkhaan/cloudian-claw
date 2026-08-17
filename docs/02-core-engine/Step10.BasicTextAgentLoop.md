# Step 10 — Basic Text Agent Loop

**Knowledge depth: 7/10**

Implement the smallest bounded agent run before adding tool execution.

## Step outcome

A request loads context, calls the provider, persists a final text response, and reports usage without shared run state.

## Task 1 — Define per-run state

### Theory

- **Architectural role:** `Runner` is a stateless orchestrator that holds ports and configuration. `runState` is one state-machine instance containing the request, temporary transcript, usage, and counters for exactly one invocation.
- **Why:** Keeping the current user, session, or messages on `Runner` creates data races and context leaks between concurrent requests. Per-run ownership also gives cancellation and accounting a clear scope.
- **GoClaw reference:** Compare request and result types in [`goclaw/internal/agent/loop_types.go`](../../goclaw/internal/agent/loop_types.go) with per-run state and substates in [`goclaw/internal/pipeline/run_state.go`](../../goclaw/internal/pipeline/run_state.go) and [`goclaw/internal/pipeline/substates.go`](../../goclaw/internal/pipeline/substates.go). The entry point is [`goclaw/internal/agent/loop_run.go`](../../goclaw/internal/agent/loop_run.go).

Read [01 — Agent Loop](../01-agent-loop.md). GoClaw describes eight logical phases:

```text
context → history → prompt → think → act → observe → memory → summarize
```

Some phases can share code, but their responsibilities must remain visible. Run state belongs to one request; a reusable runner holds dependencies only.

### Goal

Separate state that lives for one run from long-lived, shared dependencies.

### Guide to implement

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

### Theory

- **Architectural role:** This is the `context → history → prompt` phase. It resolves the agent, loads state, selects capabilities, and persists user intent before an external side effect.
- **Why:** If the provider responds before the user message is stored, the session can contain an assistant answer with no recorded cause. Phase-labeled errors identify the failure point without inventing output.
- **GoClaw reference:** Follow `Loop.Run` in [`goclaw/internal/agent/loop_run.go`](../../goclaw/internal/agent/loop_run.go), context preparation in [`goclaw/internal/agent/loop_context.go`](../../goclaw/internal/agent/loop_context.go), history assembly in [`goclaw/internal/agent/loop_history.go`](../../goclaw/internal/agent/loop_history.go), and pipeline `ContextStage` in [`goclaw/internal/pipeline/context_stage.go`](../../goclaw/internal/pipeline/context_stage.go).

### Goal

Create a valid, scoped, and durable input snapshot before the first model call.

### Guide to implement

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

- **Architectural role:** Each iteration creates one `ChatRequest`, accumulates `Usage`, appends the assistant proposal, and chooses the next transition from the canonical response. Limits turn an open-ended loop into a finite state machine.
- **Why:** The model cannot be trusted to stop by itself. Iteration, output, and context budgets are runtime resource guards and should not depend on the prompt or provider to enforce them.
- **GoClaw reference:** Inspect the production loop in [`goclaw/internal/agent/loop.go`](../../goclaw/internal/agent/loop.go) and [`goclaw/internal/agent/loop_run.go`](../../goclaw/internal/agent/loop_run.go). The easier-to-study staged version is in [`goclaw/internal/pipeline/think_stage.go`](../../goclaw/internal/pipeline/think_stage.go), [`tool_stage.go`](../../goclaw/internal/pipeline/tool_stage.go), and [`observe_stage.go`](../../goclaw/internal/pipeline/observe_stage.go).

The model either returns a final answer or proposes tool calls. The runtime executes allowed tools, appends observations, and asks the model again.

### Goal

Implement the central state transition: finish when the model returns final text, or move through act and observe before thinking again when it proposes tools.

### Guide to implement

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

Apply a temporary conservative context limit before every provider call; Step 14 implements the full budgeting and pruning policy.
