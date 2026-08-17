# Step 6 — OpenAI-Compatible Provider Basics

**Knowledge depth: 6/10**

Connect one OpenAI-compatible endpoint through the canonical provider boundary.

## Step outcome

The core can make a bounded non-streaming model call without seeing provider wire types.

## Task 1 — Define provider capabilities

### Theory

- **Architectural role:** Capabilities are adapter metadata, separate from an individual request. The agent uses them to choose streaming or non-streaming and to avoid tool or media combinations that the endpoint cannot support.
- **Why:** Sharing a `Chat` interface does not mean every backend supports every feature. Sending an unsupported request and handling the error later creates late failures and makes safe fallback harder.
- **GoClaw reference:** Inspect `ProviderCapabilities` and the adapter contract in [`goclaw/internal/providers/adapter_registry.go`](../../goclaw/internal/providers/adapter_registry.go), `OpenAIAdapter.Capabilities` in [`goclaw/internal/providers/adapter_openai.go`](../../goclaw/internal/providers/adapter_openai.go), and other adapters for comparison.

Read [02 — LLM Providers](../02-providers.md). Providers differ in streaming, tool calls, reasoning fields, and media support. [12 — Extended Thinking](../12-extended-thinking.md) explains reasoning controls. [18 — ACP Provider](../18-acp-provider.md) is reference architecture only.

Capabilities tell the agent which request path is valid before bytes are sent.

### Goal

Let the core ask “what can this provider do?” before choosing a valid request path.

### Guide to implement

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

- **Architectural role:** The translator is an anti-corruption layer: canonical requests become wire requests, and wire responses become canonical responses. Only this layer understands `choices`, function-tool envelopes, or argument JSON strings.
- **Why:** Types that look similar today may differ later in roles, media, finish reasons, or tool IDs. Explicit mapping creates one testable location for provider differences.
- **GoClaw reference:** Follow `OpenAIAdapter.ToRequest`, `FromResponse`, and `FromStreamChunk` in [`goclaw/internal/providers/adapter_openai.go`](../../goclaw/internal/providers/adapter_openai.go), then compare the canonical types in [`goclaw/internal/providers/types.go`](../../goclaw/internal/providers/types.go) with wire helpers in [`goclaw/internal/providers/openai_request.go`](../../goclaw/internal/providers/openai_request.go).

Wire types describe the provider JSON. Canonical types describe the agent. Keep them separate even when fields look similar.

### Goal

Isolate OpenAI-compatible JSON and schema details inside one two-way adapter.

### Guide to implement

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

### Theory

- **Architectural role:** The provider owns the HTTP client and protocol details. It does not own history, whole-run retries, tool execution, or session persistence.
- **Why:** `http.NewRequestWithContext` connects request lifecycle to the network call. A bounded error body prevents memory and log abuse, while an injected `http.Client` supports testing and centralized timeout configuration.
- **GoClaw reference:** Follow `OpenAIProvider.Chat` in [`goclaw/internal/providers/openai_chat.go`](../../goclaw/internal/providers/openai_chat.go), request construction in [`goclaw/internal/providers/openai_request.go`](../../goclaw/internal/providers/openai_request.go), and shared retry and error classification in [`goclaw/internal/providers/retry.go`](../../goclaw/internal/providers/retry.go) and [`goclaw/internal/providers/error_classify.go`](../../goclaw/internal/providers/error_classify.go).

### Goal

Complete the first outbound adapter with timeouts, cancellation, status handling, and safe secret handling.

### Guide to implement

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
