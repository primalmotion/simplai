package llm

import "context"

// Embedder is the embeddings interface
type Embedder interface {
	EmbedChunks(ctx context.Context, chunks []string, options ...EmbeddingOption) ([][]float64, error)
	EmbedQuery(ctx context.Context, query string, options ...EmbeddingOption) ([]float64, error)
}
