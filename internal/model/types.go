package model

import "time"

// Create custom type
// Sau do tein hanh gan cac bien moi vao trong type nay
// Value cua bien moi phai bieu dien duoc duoi dang underlying type
type SessionStatus string 
const (
	SessionActive SessionStatus = "active" 
	SessionClose SessionStatus = "close"
)
type AgentType string 
const (
	AgentPredefined AgentType = "predefined" 
	AgentOpen AgentType = "open"
)
// import "context"
type Agent struct {
	ID string `json:"id"`
	Name string `json:"name"` //agent name 
	Kind string `json:"kind"` //Agent soul: coder, artist, teacher... 
	SystemPrompt string `json:"system_prompt"`
	Type AgentType `json:"type,omitempty"`
	MaxToken int `json:"max_token"`
	//Mot so chi so config khac, se tien hanh nghien cuu lai database va import sau 
	ProviderModel ProviderModel `json:"provider_model"`
}

type ImageContent struct {
	MimeType string `json:"mime_type"`
	Data string `json:"data"` //base64 encoded image bytes 
	URL string `json:"url,omitempty"`
	Size int `json:"size,omitempty"`
} 

type VideoContent struct {
	MimeType string `json:"mime_type"`
	Data string `json:"data"`// base64-encoded video bytes
	URL      string `json:"url,omitempty"`     // URL of the video
	Partial  bool   `json:"partial,omitempty"`
}
//Can phai tao them 1 lop tham chieu den media data de truyen vao trong message 
type MediaRef struct {
	ID string `json:"id"`
	MimeType string `json:"mime_type"`
	Kind string `json:"kind"` 
	URL string `json:"url,omitempty"`
}
//Message represent as a conversation message to an AI model 
type Message struct {
	Role string `json:"role"` //system, assistant, user, tool
	Content string `json:"content"` //
	Thinking string `json:"thinking,omitempty"` //Reasoning content for thinking models (Kimi, Deepseek, etc...)
	Images []ImageContent `json:"-"` //Khong serialize 2 file nay vi rat lon khi chuyen qua json 
	Videos []VideoContent `json:"-"`
	MediaRefs []MediaRef `json:"media_refs,omitempty"`
	Tools []Tool `json:"tools,omitempty"`
	Rules []Rule `json:"rules,omitempty"`
}
//LLM provider -> It will include base url so that you can call to it 
type Provider struct {
	ID string `json:"id"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`   
	BaseURL string `json:"base_url"`
	APIKeys []APIKey  
	//Khng can thiet phai dat het relationship, chi can khai bao domain don gian la duoc. Thieu thi se tien hanh bo sung sau
}


type Model struct {
	ID string `json:"id"` //unique id of model 
	Name string `json:"name,omitempty"` //model name, don't depends on the provider 
	Description string `json:"description,omitempty"`
}

type ProviderModel struct {
	ProviderId string `json:"provider_id"`
	ModelId string `json:"model_id"`
	ProviderConfigId string `json:"provider_config_id"` //Id for LLM model that provider quy định 
}

type APIKey struct {
	Value string `json:"value"`
	IsUsed bool 
	IsExpired bool 
}
//Tien hanh phan tach ra thanh nhieu dang interface: tool_definitions, tool_calls, tool_function_schema theo kien truc cua goclaw 

type Tool struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Arguments map[string]any `json:"arguments"`
}

type Skill struct {
	Name string `json:"name"` 
	Summary string `json:"summary"`
	Content string `json:"content,omitempty"`
	Enabled bool `json:"enabled"`
}

type Rule struct {
	Name string `json:"name"` 
	Summary string `json:"summary"`
	Content string `json:"content,omitempty"`
	Enabled bool `json:"enabled"`
}

type Session struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Agent Agent `json:"agent"`
	Status SessionStatus `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` 
}
//Usage per request 
type Usage struct {
	PromptTokens                      int  `json:"prompt_tokens"`
	CompletionTokens                  int  `json:"completion_tokens"`
	TotalTokens                       int  `json:"total_tokens"`
	CacheCreationTokens               int  `json:"cache_creation_input_tokens,omitempty"`
	CacheReadTokens                   int  `json:"cache_read_input_tokens,omitempty"`
	PromptTokensIncludeCachedSegments bool `json:"prompt_tokens_include_cached_segments,omitempty"`
	ThinkingTokens                    int  `json:"thinking_tokens,omitempty"`
	RequestCount                      int  `json:"request_count,omitempty"`
	ImageCount                        int  `json:"image_count,omitempty"`
	WebSearchCount                    int  `json:"web_search_count,omitempty"`
}

type Memory struct {
	ID string `json:"id"`
	Content string `json:"content"`
	Type MemoryType `json:"type"`
	Agent Agent `json:"agent"`
	AgentID string `json:"agent_id,omitempty"`
}
type MemoryType string 
const (
	MemorySearchTool MemoryType = "memory_search_tool"
	MemorySearchSkill MemoryType = "memory_search_skill"
)

type Execution struct {
	
}