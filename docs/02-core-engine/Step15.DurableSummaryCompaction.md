# Step 15 — Durable Summary Compaction

**Knowledge depth: 8/10**

Compact old coherent turns into a durable, recoverable summary checkpoint.

## Step outcome

Repeated compaction survives restart without losing recent messages or summarizing the same range twice.

## Task 1 — Save a durable summary checkpoint

### Theory

- **Architectural role:** The compactor selects a coherent range, the provider creates a factual summary, and the store commits the summary with a watermark or range. The prompt is then rebuilt from the checkpoint and remaining suffix.
- **Why:** A summary kept only in memory disappears after restart. Removing source messages before saving the checkpoint risks data loss, while omitting a watermark can summarize the same messages repeatedly.
- **GoClaw reference:** Inspect compaction orchestration in [`goclaw/internal/agent/loop_compact.go`](../../goclaw/internal/agent/loop_compact.go), pipeline checkpoints in [`goclaw/internal/pipeline/checkpoint_stage.go`](../../goclaw/internal/pipeline/checkpoint_stage.go), and `GetSummary` and `SetSummary` persistence in [`goclaw/internal/store/pg/sessions.go`](../../goclaw/internal/store/pg/sessions.go).

A summary sent only to the next model call is unsafe. A process crash would lose the compacted context.

### Goal

Compact an old prefix into a durable checkpoint that can resume after a crash or restart.

### Guide to implement

When compaction is required:

1. Select the oldest coherent message range.
2. Ask the provider for a factual, instruction-free summary.
3. Store the new summary.
4. Mark or delete only the summarized range in the same transaction or a recoverable sequence.
5. Rebuild the prompt and confirm it fits the budget.

Include durable facts, decisions, unresolved questions, and important references. Exclude hidden reasoning and any claim not present in the source messages.

## Task 2 — Verify long-session behavior

### Theory

- **Architectural role:** The fixture must contain text, tool groups, and an oversized result so it tests the full budgeting pipeline rather than only one counting helper.
- **Why:** Context bugs often appear only during a second compaction or after reload. Idempotency and watermark tests catch summaries of summaries and repeated processing of the same prefix.
- **GoClaw reference:** Review [`goclaw/internal/agent/loop_compact_integration_test.go`](../../goclaw/internal/agent/loop_compact_integration_test.go), [`goclaw/internal/pipeline/compaction_pressure_e2e_test.go`](../../goclaw/internal/pipeline/compaction_pressure_e2e_test.go), and preservation tests in [`goclaw/internal/pipeline/compaction_pending_preservation_test.go`](../../goclaw/internal/pipeline/compaction_pending_preservation_test.go).

### Goal

Prove that the invariants still hold across repeated pruning, compaction, and restarts.

### Guide to implement

Create a fixture with many turns, at least two tool-call pairs, and one oversized tool result. Assert that:

- The final prompt is within budget.
- Recent turns remain exact.
- Tool pairs remain complete and ordered.
- Truncation is visible.
- The summary survives a new store instance.
- Repeated compaction does not summarize the same range twice.

This step is complete when a long session can continue after restart without invalid provider history.
