package entity

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