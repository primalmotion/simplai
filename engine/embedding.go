package engine

import "context"

// Embedder is the embeddings interface
type Embedder interface {
	Embed(ctx context.Context, chunks []string, options ...EmbeddingOption) ([][]float64, error)
}
