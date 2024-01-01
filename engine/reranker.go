package engine

import "context"

// Reranker is the main interface to interact with a Reranker.
type Reranker interface {
	Rerank(ctx context.Context, query string, passages []string) ([]float64, error)
}
