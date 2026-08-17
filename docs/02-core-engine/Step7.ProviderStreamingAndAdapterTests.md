# Step 7 — Provider Streaming and Adapter Tests

**Knowledge depth: 7/10**

Add safe SSE streaming and verify the complete provider adapter boundary.

## Step outcome

Text and fragmented tool calls stream correctly, cancellation works, and secrets stay out of errors.

## Task 1 — Implement streaming

### Theory

- **Architectural role:** The SSE parser handles framing, an accumulator joins text and tool-argument fragments, and a callback sends progress toward the gateway. Persistence and the decision to execute a tool remain in the agent layer.
- **Why:** Tool arguments may be split at any byte or frame boundary. Retrying after a delta has been emitted breaks at-most-once visibility and can repeat a tool proposal.
- **GoClaw reference:** Study `ChatStream` and its accumulator in [`goclaw/internal/providers/openai_chat.go`](../../goclaw/internal/providers/openai_chat.go), the shared scanner in [`goclaw/internal/providers/sse_reader.go`](../../goclaw/internal/providers/sse_reader.go), and fragmentation fixtures in [`goclaw/internal/providers/sse_reader_test.go`](../../goclaw/internal/providers/sse_reader_test.go).

OpenAI-compatible streaming normally uses server-sent events. Tool arguments may arrive as fragments and must be joined before JSON decoding.

### Goal

Turn a long-running response into observable deltas while still producing one complete canonical `ChatResponse`.

### Guide to implement

In `ChatStream`:

1. Send the same request with `stream: true`.
2. Scan SSE `data:` frames with a configured maximum frame size.
3. Emit text deltas through `onChunk`.
4. Accumulate tool-call fragments by choice index and tool-call index.
5. On `[DONE]`, assemble one canonical `ChatResponse`.
6. Stop immediately when the callback or context returns an error.

Do not automatically retry after a text delta has reached the caller. A retry could duplicate visible output or repeat a tool proposal.

## Task 2 — Test the adapter boundary

### Theory

- **Architectural role:** The test uses the adapter's public boundary: it sends a canonical request, observes the HTTP wire format, returns a wire fixture, and asserts the canonical response or error.
- **Why:** Mocking `Chat` directly tests only the mock. `httptest.Server` can verify real headers, URLs, bodies, SSE framing, cancellation, and response limits.
- **GoClaw reference:** Review test patterns in [`goclaw/internal/providers/openai_test.go`](../../goclaw/internal/providers/openai_test.go), request contracts in [`goclaw/internal/providers/openai_request_test.go`](../../goclaw/internal/providers/openai_request_test.go), and streaming edge cases in [`goclaw/internal/providers/sse_reader_ctx_test.go`](../../goclaw/internal/providers/sse_reader_ctx_test.go).

### Goal

Lock down the translation and network contract with a local test server instead of a real external API.

### Guide to implement

Use `httptest.Server` to cover:

1. A normal assistant text response.
2. One tool call with valid JSON arguments.
3. A streamed response split across several SSE frames.
4. Fragmented tool arguments.
5. A non-2xx response with a large body.
6. Context cancellation.

Add a manual smoke command that sends a short prompt to the configured endpoint. The step is complete when the core receives only canonical responses and logs contain no provider secret.
