package entity

import (
	"errors"
)

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

func (u Usage) Validate() error {
	if u.PromptTokens < 0 {
		return errors.New("validate usage: prompt_tokens must not be negative")
	}
	if u.CompletionTokens < 0 {
		return errors.New("validate usage: completion_tokens must not be negative")
	}
	if u.TotalTokens < 0 {
		return errors.New("validate usage: total_tokens must not be negative")
	}
	if u.CacheCreationTokens < 0 {
		return errors.New("validate usage: cache_creation_tokens must not be negative")
	}
	if u.CacheReadTokens < 0 {
		return errors.New("validate usage: cache_read_tokens must not be negative")
	}
	if u.ThinkingTokens < 0 {
		return errors.New("validate usage: thinking_tokens must not be negative")
	}
	if u.RequestCount < 0 {
		return errors.New("validate usage: request_count must not be negative")
	}
	if u.ImageCount < 0 {
		return errors.New("validate usage: image_count must not be negative")
	}
	if u.WebSearchCount < 0 {
		return errors.New("validate usage: web_search_count must not be negative")
	}
	return nil
}
