package entity
type Embedding []float64 
type MemoryDocument struct {
	ID string `json:"id"`
	AgentID string `json:"agent_id"`
	Agent Agent 
	Workspace string `json:"workspace,omitempty"`
	Path string `json:"path,omitempty"`
	FileName string `json:"file_name,omitempty"`
	Content string `json:"content"`
	Hash string `json:"hash"`
}

type MemoryTrunk struct {
	ID string `json:"id"`
	DocumentID string `json:"document_id"` 
	StartLine int `json:"start_line"`
	EndLine int `json:"end_line"`
	Content string `json:"content"`
	Embedding Embedding 
	Hash string `json:"hash"`
}