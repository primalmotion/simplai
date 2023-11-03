package embeddings

import "context"

// Embedder is the embeddings interface
type Embedder interface {
	EmbedChunks(ctx context.Context, chunks []string) ([][]float64, error)
	EmbedQuery(ctx context.Context, query string) ([]float64, error)
}
