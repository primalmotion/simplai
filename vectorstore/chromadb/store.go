package chromadb

import (
	"context"
	"fmt"

	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/vectorstore"
)

var _ vectorstore.VectorStore = &ChromaStore{}

// A ChromaStore implement the the
// vectorstore.VectorStore interface for a
// a store backend by ChromaDB.
type ChromaStore struct {
	embedder     engine.Embedder
	client       *Client
	collectionID string
}

// NewChromaStore returns a new *ChromaStore.
func NewChromaStore(client *Client, collectionID string, embedder engine.Embedder) *ChromaStore {
	return &ChromaStore{
		client:       client,
		embedder:     embedder,
		collectionID: collectionID,
	}
}

// AddDocument implements the vectorstore.VectorStore interface.
// if perform a chromadb upsert.
func (c *ChromaStore) AddDocument(ctx context.Context, documents ...vectorstore.Document) (err error) {

	l := len(documents)
	embeddings := make([]vectorstore.Embedding, l)
	metadatas := make([]vectorstore.Metadata, l)
	contents := make([]string, l)
	ids := make([]string, l)

	for i, d := range documents {
		embeddings[i] = d.Embedding
		if len(d.Embedding) == 0 {
			embeddings[i], err = c.embedder.EmbedQuery(ctx, d.Content)
			if err != nil {
				return fmt.Errorf("unable to embedd document: %w", err)
			}
		}
		contents[i] = d.Content
		ids[i] = d.ID
		metadatas[i] = d.Metadata
	}

	err = c.client.Upsert(
		ctx,
		c.collectionID,
		EmbeddingUpdate{
			Embeddings: embeddings,
			Metadatas:  metadatas,
			Documents:  contents,
			IDs:        ids,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to execute upsert: %w", err)
	}

	return nil
}

// SimilaritySearch implements the vectorstore.VectorStore interface.
// if perform a chromadb query.
func (c *ChromaStore) SimilaritySearch(ctx context.Context, input vectorstore.Embedding, max int) ([]vectorstore.Document, error) {

	res, err := c.client.Query(
		ctx,
		c.collectionID,
		EmbeddingQuery{
			QueryEmbeddings: []vectorstore.Embedding{input},
			NResults:        max,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}

	out := make([]vectorstore.Document, len(res.Distances))
	for i := range res.Distances {
		out[i] = vectorstore.Document{
			ID:       res.IDs[0][i],
			Content:  res.Documents[0][i],
			Distance: res.Distances[0][i],
			Metadata: res.Metadatas[0][i],
		}
	}

	return out, nil
}
