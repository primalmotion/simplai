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

type VectorStore interface {
	AddDocument(context.Context, ...Document) error
	SimilaritySearch(context.Context, Embedding, int) ([]Document, error)
}
