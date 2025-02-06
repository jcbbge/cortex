package memory

import (
	"context"
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMClient defines the interface for LLM operations
type LLMClient interface {
	// GetEmbedding gets a vector embedding for text
	GetEmbedding(ctx context.Context, text string) ([]float32, error)
	
	// Chat sends a series of messages and returns streaming response
	Chat(ctx context.Context, messages []Message) (chan string, chan uint64, chan error)
}