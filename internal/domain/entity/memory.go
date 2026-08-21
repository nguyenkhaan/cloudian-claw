package entity

import (
	"errors"
)

type Embedding []float64

type MemoryDocument struct {
	ID        string `json:"id"`
	AgentID   string `json:"agent_id"`
	Agent     Agent
	Workspace string `json:"workspace,omitempty"`
	Path      string `json:"path,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	Content   string `json:"content"`
	Hash      string `json:"hash"`
}

type MemoryTrunk struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	StartLine  int    `json:"start_line"`
	EndLine    int    `json:"end_line"`
	Content    string `json:"content"`
	Embedding  Embedding
	Hash       string `json:"hash"`
}

func (d MemoryDocument) Validate() error {
	if d.ID == "" {
		return errors.New("validate memory_document: id is required")
	}
	if d.AgentID == "" {
		return errors.New("validate memory_document: agent_id is required")
	}
	if d.Content == "" {
		return errors.New("validate memory_document: content is required")
	}
	if d.Hash == "" {
		return errors.New("validate memory_document: hash is required")
	}
	return nil
}

func (t MemoryTrunk) Validate() error {
	if t.ID == "" {
		return errors.New("validate memory_trunk: id is required")
	}
	if t.DocumentID == "" {
		return errors.New("validate memory_trunk: document_id is required")
	}
	if t.StartLine > t.EndLine {
		return errors.New("validate memory_trunk: start_line must not be greater than end_line")
	}
	return nil
}
