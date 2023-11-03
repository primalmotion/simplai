package ollama

import (
	"context"
	"errors"

	"github.com/primalmotion/simplai/embeddings"
	"github.com/primalmotion/simplai/llm"
	ollamaclient "github.com/primalmotion/simplai/llm/ollama/internal"
)

// EmbedChunks implement the embeddings interface for chunks.
func (e *ollamaAPI) EmbedChunks(ctx context.Context, chunks []string, options ...embeddings.Option) ([][]float64, error) {

	opts := defaultEmbeddingConfig()
	for _, opt := range options {
		opt(&opts)
	}

	emb := make([][]float64, 0, len(chunks))

	batches := embeddings.Batch(chunks, opts.BatchSize)
	for _, batch := range batches {

		currentEmbeddings := [][]float64{}

		for _, chunk := range chunks {
			embedding, err := e.client.Embed(ctx, &ollamaclient.EmbeddingRequest{
				Prompt: chunk,
				Model:  opts.Model,
			})
			if err != nil {
				return nil, err
			}

			if len(embedding.Embedding) == 0 {
				return nil, errors.New("no response")
			}

			currentEmbeddings = append(currentEmbeddings, embedding.Embedding)
		}

		if len(chunks) != len(currentEmbeddings) {
			return currentEmbeddings, errors.New("no all input got emmbedded")
		}

		// get num of token in that batch
		// we should use the encoder of the model to get the tokens
		// but its not available. So we fall back on tiktoken
		numTokens := make([]float64, 0, len(batch))
		for _, text := range batch {
			numTokens = append(numTokens, float64(llm.CountTokens(opts.Model, text)))
		}

		if len(currentEmbeddings) > 1 {
			combinedVectors, err := embeddings.CombineBatchedEmbedding(currentEmbeddings, numTokens)
			if err != nil {
				return [][]float64{}, err
			}
			emb = append(emb, combinedVectors)
			continue
		}

		emb = append(emb, currentEmbeddings...)
	}

	return emb, nil
}

// EmbedQuery implement the embeddings interface for query.
func (e *ollamaAPI) EmbedQuery(ctx context.Context, query string) ([]float64, error) {
	c, err := e.EmbedChunks(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	return c[0], nil
}
