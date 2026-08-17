# Step 4 — LLM Provider Adapter

**Knowledge depth: 6/10**

This step connects one OpenAI-compatible model endpoint without leaking provider wire types into the agent core.

## Task 1 — Define provider capabilities

### Theory

Read [02 — LLM Providers](../02-providers.md). Providers differ in streaming, tool calls, reasoning fields, and media support. [12 — Extended Thinking](../12-extended-thinking.md) explains reasoning controls. [18 — ACP Provider](../18-acp-provider.md) is reference architecture only.

Capabilities tell the agent which request path is valid before bytes are sent.

### Practice guide

Add these types to `internal/model`:

```go
type Capabilities struct {
	Streaming       bool
	ToolCalling     bool
	StreamWithTools bool
}

type StreamChunk struct {
	Text string
	Done bool
}

type StreamingProvider interface {
	Provider
	Capabilities() Capabilities
	ChatStream(
		ctx context.Context,
		req ChatRequest,
		onChunk func(StreamChunk) error,
	) (ChatResponse, error)
}
```

For this project, report streaming and tool calling as supported. If the selected endpoint cannot stream tool calls, set `StreamWithTools` to false and use non-streaming whenever tools are present.

## Task 2 — Create wire types and translators

### Theory

Wire types describe the provider JSON. Canonical types describe the agent. Keep them separate even when fields look similar.

### Practice guide

Create `internal/providers/openai/types.go` and `mapping.go`. Implement:

```go
func toOpenAIRequest(model.ChatRequest, bool) openAIRequest
func fromOpenAIResponse(openAIResponse) (model.ChatResponse, error)
```

Map:

- Canonical roles and content.
- Tool definitions to OpenAI function tools.
- Assistant tool calls and matching tool results.
- Usage and finish reason.
- JSON tool arguments into `map[string]any`.

If a tool call has no ID, generate one. If arguments are malformed, return a model-visible tool error later; do not guess values.

## Task 3 — Implement non-streaming chat

### Practice guide

Create `internal/providers/openai/provider.go`:

```go
type Provider struct {
	baseURL     string
	apiKey      string
	defaultModel string
	httpClient  *http.Client
}

func (p *Provider) Chat(ctx context.Context, req model.ChatRequest) (model.ChatResponse, error) {
	body, err := json.Marshal(toOpenAIRequest(req, false))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("encode provider request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimRight(p.baseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("create provider request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("call provider: %w", err)
	}
	defer resp.Body.Close()
	// Check status, decode a bounded body, then call fromOpenAIResponse.
}
```

Use a client timeout and also honor `ctx`. Limit non-success response bodies, for example to 64 KiB. Never include the API key in an error.

## Task 4 — Implement streaming

### Theory

OpenAI-compatible streaming normally uses server-sent events. Tool arguments may arrive as fragments and must be joined before JSON decoding.

### Practice guide

In `ChatStream`:

1. Send the same request with `stream: true`.
2. Scan SSE `data:` frames with a configured maximum frame size.
3. Emit text deltas through `onChunk`.
4. Accumulate tool-call fragments by choice index and tool-call index.
5. On `[DONE]`, assemble one canonical `ChatResponse`.
6. Stop immediately when the callback or context returns an error.

Do not automatically retry after a text delta has reached the caller. A retry could duplicate visible output or repeat a tool proposal.

## Task 5 — Test the adapter boundary

### Practice guide

Use `httptest.Server` to cover:

1. A normal assistant text response.
2. One tool call with valid JSON arguments.
3. A streamed response split across several SSE frames.
4. Fragmented tool arguments.
5. A non-2xx response with a large body.
6. Context cancellation.

Add a manual smoke command that sends a short prompt to the configured endpoint. The step is complete when the core receives only canonical responses and logs contain no provider secret.
