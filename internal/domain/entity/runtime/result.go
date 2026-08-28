package runtime

import "cloudian/cloudian-claw/internal/domain/entity"

type ModelResponse struct {
	ToolCalls []ToolCall 
	Output string 
	Model string 
	TotalCall int //So lan da goi 
	Usage entity.Usage
	ErrorMessage string //Neu nhu LLM tra loi thi loi cu the la gi 
} 

type ExecutionResult struct {
	ExecutionID string 
	Output string 
	ErrorMessage string 
	ToolCalls []ToolCall 
	Usage entity.Usage 
	Status ExecutionStatus 
}