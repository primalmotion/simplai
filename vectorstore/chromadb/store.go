package chromadb

import (
	"context"
	"fmt"

	"github.com/primalmotion/simplai/vectorstore"
)

var _ vectorstore.VectorStore = &ChromaStore{}

type ChromaStore struct {
	client       *Client
	collectionID string
}

func NewChromaStore(client *Client, collectionID string) *ChromaStore {
	return &ChromaStore{
		client:       client,
		collectionID: collectionID,
	}
}

func (c *ChromaStore) AddDocument(ctx context.Context, documents ...vectorstore.Document) error {

	l := len(documents)
	embeddings := make([]vectorstore.Embedding, l)
	metadatas := make([]vectorstore.Metadata, l)
	contents := make([]string, l)
	ids := make([]string, l)

	for i, d := range documents {
		embeddings[i] = d.Embedding
		contents[i] = d.Content
		ids[i] = d.ID
		metadatas[i] = d.Metadata
	}

	c.client.Upsert(
		ctx,
		c.collectionID,
		EmbeddingUpdate{
			Embeddings: embeddings,
			Metadatas:  metadatas,
			Documents:  contents,
			IDs:        ids,
		},
	)

	return nil
}

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
