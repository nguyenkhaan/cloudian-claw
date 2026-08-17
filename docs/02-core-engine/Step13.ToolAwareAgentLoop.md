# Step 13 — Tool-Aware Agent Loop

**Knowledge depth: 8/10**

Extend the basic agent run with act-observe iterations and explicit failure semantics.

## Step outcome

The runner completes text-only and tool round trips while preserving call IDs, message order, budgets, and cancellation.

## Task 1 — Execute and observe tool calls

### Theory

- **Architectural role:** The agent loop controls order and persistence, the registry validates, authorizes, and executes, and each result is correlated by call ID. An assistant tool request and its tool results form one inseparable group.
- **Why:** The model has no execution authority. Persisting the proposal before the action improves audit and recovery, while preserving original order keeps the next provider request valid even if read-only work runs concurrently.
- **GoClaw reference:** Follow the tool path in [`goclaw/internal/agent/loop_tools.go`](../../goclaw/internal/agent/loop_tools.go), pipeline callbacks in [`goclaw/internal/agent/loop_pipeline_tool_callbacks.go`](../../goclaw/internal/agent/loop_pipeline_tool_callbacks.go), and registry execution in [`goclaw/internal/tools/registry.go`](../../goclaw/internal/tools/registry.go).

The model proposes actions. The tool registry independently validates and authorizes them. Unknown, invalid, and denied calls become canonical tool error messages so the model may recover.

### Goal

Turn a model proposal into a runtime-controlled observation and return it to the transcript in a provider-valid form.

### Guide to implement

For every tool-call batch:

1. Check the remaining total tool-call budget.
2. Persist the assistant message containing the original calls.
3. Execute each call through the registry.
4. Append and persist one tool message for each call ID.
5. Preserve provider order even if read-only I/O runs concurrently.
6. Continue to the next model iteration.

Never keep an assistant tool request without its matching tool results. Never execute the same call twice because persistence returned an ambiguous error.

## Task 2 — Define failure and cancellation behavior

### Theory

- **Architectural role:** Cancellation propagates through providers and tools, recoverable tool errors become observations, infrastructure or run-limit errors stop the run, and final persistence is the commit point required before reporting completion.
- **Why:** Converting every failure into assistant text can make users believe data was saved or an action ran successfully. `WithoutCancel` is suitable only for bounded cleanup, not for continuing business work after cancellation.
- **GoClaw reference:** Review cancellation tests in [`goclaw/internal/agent/loop_cancel_test.go`](../../goclaw/internal/agent/loop_cancel_test.go), finalization in [`goclaw/internal/agent/loop_finalize.go`](../../goclaw/internal/agent/loop_finalize.go), pipeline finalization in [`goclaw/internal/pipeline/finalize_stage.go`](../../goclaw/internal/pipeline/finalize_stage.go), and error paths in [`goclaw/internal/agent/loop_run.go`](../../goclaw/internal/agent/loop_run.go).

Read the run-visibility concepts in [10 — Tracing & Observability](../10-tracing-observability.md). A clear failure is better than a fluent but false assistant message.

### Goal

Define terminal states and durability rules for each failure type before they occur in production.

### Guide to implement

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

## Task 3 — Test the complete loop

### Theory

- **Architectural role:** A fake provider acts as a deterministic script, a fake store records operations, and fake tools create observations. Assertions focus on order, correlation, budgets, and isolation.
- **Why:** A real model is non-deterministic, slow, and costly for unit tests. The concurrency test specifically proves that `Runner` does not keep mutable per-run state.
- **GoClaw reference:** Review [`goclaw/internal/agent/toolloop_test.go`](../../goclaw/internal/agent/toolloop_test.go), [`goclaw/internal/agent/loop_cancel_test.go`](../../goclaw/internal/agent/loop_cancel_test.go), and stage tests in [`goclaw/internal/pipeline/pipeline_test.go`](../../goclaw/internal/pipeline/pipeline_test.go).

### Goal

Test the state machine with scripted responses that deliberately cover each transition and failure point.

### Guide to implement

Use scripted fake provider responses to test:

1. One text-only response.
2. One tool call followed by a final response.
3. Unknown tool followed by model recovery.
4. Repeated calls stopped by the tool-call limit.
5. Provider cancellation.
6. Persistence failure at each append point.

Assert message order, tool-call IDs, usage totals, and iteration counts. This step is complete when one runner can handle concurrent calls without sharing run state.
