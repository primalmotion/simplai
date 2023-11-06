package memdb

import (
	"context"
	"fmt"
	"sort"

	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/vectorstore"
	"gonum.org/v1/gonum/mat"
)

// A MemoryStore is an in memory simple vector db.
type MemoryStore struct {
	embedder engine.Embedder
	db       map[string]vectorstore.Document
}

// New return a new memory store.
func New(embedder engine.Embedder) *MemoryStore {
	return &MemoryStore{
		embedder: embedder,
		db:       make(map[string]vectorstore.Document),
	}
}

// AddDocument implement the vectorstore interface.
func (m *MemoryStore) AddDocument(ctx context.Context, documents ...vectorstore.Document) error {
	for _, d := range documents {
		if len(d.Embedding) == 0 {
			em, err := m.embedder.Embed(ctx, []string{d.Content})
			if err != nil {
				return fmt.Errorf("unable to embedd document: %w", err)
			}
			d.Embedding = em[0]
		}
		m.db[d.ID] = d
	}

	return nil
}

// SimilaritySearch implement the vectorstore interface.
func (m *MemoryStore) SimilaritySearch(ctx context.Context, query string, max int) ([]vectorstore.Document, error) {

	queryEmbedding, err := m.embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, err
	}

	scores := make(map[float64]string, len(m.db))
	sl := make([]float64, 0, len(m.db))
	for id, d := range m.db {
		score := cosineSimilarity(queryEmbedding[0], d.Embedding)
		scores[score] = id
		sl = append(sl, score)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(sl)))

	count := max
	if len(sl) < max {
		count = len(sl)
	}

	res := make([]vectorstore.Document, 0, count)

	for _, dist := range sl[:count] {
		pd := m.db[scores[dist]]
		r := vectorstore.Document{
			Metadata:  pd.Metadata,
			ID:        pd.ID,
			Content:   pd.Content,
			Embedding: pd.Embedding,
			Distance:  dist,
		}
		res = append(res, r)
	}

	return res, nil
}

// cosineSimilarityi permform the cosine similarity between two embeddings.
// this is the v1.v2/(||v1||*||v2||)
func cosineSimilarity(e1, e2 []float64) float64 {

	v1 := mat.NewVecDense(len(e1), e1)
	v2 := mat.NewVecDense(len(e2), e2)
	normEpsilon := 1e-30

	return mat.Dot(v1, v2) / ((mat.Norm(v1, 2) + normEpsilon) * (mat.Norm(v2, 2) + normEpsilon))
}
