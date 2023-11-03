package ollama

import (
	"context"

	"github.com/primalmotion/simplai/embeddings"
	"github.com/primalmotion/simplai/llm"
	"github.com/primalmotion/simplai/llm/ollama"
)

type ollamaEmbedder struct {
	client *ollama.Client
	config embeddings.EmbeddingConfig
}

func defaultEmbeddingConfig() embeddings.EmbeddingConfig {
	return embeddings.EmbeddingConfig{
		BatchSize: 512,
	}
}

// New returns a new ollamaEmbedder
func New(c *ollama.Client, opts ...embeddings.Option) (*ollamaEmbedder, error) { //nolint:revive

	o := defaultEmbeddingConfig()
	for _, opt := range opts {
		opt(&o)
	}

	return &ollamaEmbedder{
		client: c,
		config: o,
	}, nil
}

// EmbedChunks implement the embeddings interface for chunks.
func (e *ollamaEmbedder) EmbedChunks(ctx context.Context, chunks []string) ([][]float64, error) {

	emb := make([][]float64, 0, len(chunks))

	batches := embeddings.Batch(chunks, e.config.BatchSize)
	for _, batch := range batches {

		cemb, err := e.client.Embed(ctx, batch, llm.OptionModel(e.config.Model))
		if err != nil {
			return nil, err
		}

		// get num of token in that batch
		// we should use the encoder of the model to get the tokens
		// but its not available. So we fall back on tiktoken
		numTokens := make([]float64, 0, len(batch))
		for _, text := range batch {
			numTokens = append(numTokens, float64(llm.CountTokens(e.config.Model, text)))
		}

		combinedVectors, err := embeddings.CombineBatchedEmbedding(cemb, numTokens)
		if err != nil {
			return [][]float64{}, err
		}

		emb = append(emb, combinedVectors)
	}

	return emb, nil
}

// EmbedQuery implement the embeddings interface for query.
func (e *ollamaEmbedder) EmbedQuery(ctx context.Context, query string) ([]float64, error) {
	c, err := e.EmbedChunks(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	return c[0], nil
}
