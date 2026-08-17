package model

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"` //Content of the message 
	ToolCallID string `json:"tool_call_id,omitempty"`
	ToolName string `json:"tool_name,omitempty"`
	ToolCalls []ToolCall //What is ToolCall? 
	IsError bool //Is the message have error? 
} 

//Describe a tool calling (function calling) that LLM will send to BE 
type ToolCall struct {
	ID string `json:"id"`
	Name string `json:"name,omitempty"`
	Arguments map[string]any `json:"arguments,omitempty"`
}
//Description of a tool yeahhhh 
type ToolDefinition struct {
	Name string `json:"name,omitempty"` 
	Description string `json:"description,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
}
//Tracking the total token usage per user's request 
type Usage struct {
	PromptTokens, CompletionTokens, TotalTokens int 
}
//User request 
type ChatRequest struct {
	Messages []Message        `json:"messages"`
	Tools    []ToolDefinition `json:"tools,omitempty"`
	Model    string           `json:"model,omitempty"`
	Options  map[string]any   `json:"options,omitempty"`
}
//Agent response 
type ChatResponse struct {
	Content      string     `json:"content"`
	Thinking     string     `json:"thinking,omitempty"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	FinishReason string     `json:"finish_reason"` // "stop", "tool_calls", "length"
	Usage        *Usage     `json:"usage,omitempty"`
	//We will provide those field later 

	// // Phase is Codex-specific (gpt-5.3-codex): "commentary" or "final_answer".
	// // Agent loop must persist this on assistant messages for Codex performance.
	// Phase string `json:"phase,omitempty"`

	// // RawAssistantContent preserves the raw content blocks array from the provider response.
	// // Used by Anthropic to pass thinking blocks back in tool use loops (required by API).
	// RawAssistantContent json.RawMessage `json:"-"`

	// // ThinkingSignature is the accumulated signature from streaming thinking blocks.
	// // Required by Anthropic API for tool use passback when thinking is enabled.
	// ThinkingSignature string `json:"-"`

	// // Images holds generated images returned by image_generation_call tools (Codex).
	// // Not persisted to DB; populated at runtime from provider response.
	// Images []ImageContent `json:"-"`
}
/*
Đầu tiên phải hiểu được business, các thành phần, chức năng chúng ta sẽ tiến hành làm trong dự án 
Phân rã các thành phần, chức năng thành nhiều module. Mỗi module gồm mục tiêu, yêu cầu triển khai, hướng dẫn cơ bản 
Đọc từng module và làm theo, chỉnh sửa theo kiến trúc 

Hoặc nếu bạn biết rồi thì có thể đối xử AI như một thợ code thay 
*/ 