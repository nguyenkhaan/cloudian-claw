package runtime 

type ToolCall struct {
	ToolID string 
	ToolName string 
	ID string 
	Arguments map[string]any  
}

type ToolResult struct {
	ToolCallID string 
	Output string 
} 
type ToolError struct {
	ToolCallID string 
	Code string 
	Error string 
}