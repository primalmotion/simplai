package vectorstore

import (
	"context"
	"fmt"
)

// Embedding holds chromadb embeddings.
// This only supports float64 for now.
type Embedding []float64

// Metadata holds chromadb metadata.
type Metadata map[string]any

// A Document represents an embedded document.
type Document struct {
	Metadata  Metadata
	ID        string
	Content   string
	Embedding Embedding
	Distance  float64
}

func (d Document) String() string {
	return fmt.Sprintf("<document id:%s distance:%f>", d.ID, d.Distance)
}

// A VectorStore is the interface that must implement
// all vector databases.
type VectorStore interface {
	AddDocument(context.Context, ...Document) error
	SimilaritySearch(context.Context, Embedding, int) ([]Document, error)
}
