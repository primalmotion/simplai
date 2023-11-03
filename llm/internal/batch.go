package tools

import (
	"errors"

	"gonum.org/v1/gonum/floats"
)

// Batch will split inputs by the batch size.
func Batch(inputs []string, batch int) [][]string {
	batches := make([][]string, len(inputs))

	for i, input := range inputs {
		text := []rune(input)
		for j := 0; j < len(text); j += batch {
			if j+batch >= len(text) {
				batches[i] = append(batches[i], string(text[j:]))
				break
			}

			batches[i] = append(batches[i], string(text[j:j+batch]))
		}
	}

	return batches
}

// CombineBatchedEmbedding combine the batched results into a normalized
// vector.
func CombineBatchedEmbedding(vectors [][]float64, weights []float64) ([]float64, error) {

	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, errors.New("vectors must not be empty")
	}

	if len(vectors) != len(weights) {
		return nil, errors.New("length of weights must match the number of vectors")
	}

	if len(vectors[0]) != len(vectors[1]) {
		return nil, errors.New("vectors must have the same dimension")
	}

	// Compute the weighted average
	average := make([]float64, len(vectors[0]))
	for i, vector := range vectors {
		floats.AddScaled(average, weights[i], vector)
	}

	// Calculate the norm of the average vector
	norm := floats.Norm(average, 2)

	if norm == 0 {
		return nil, errors.New("the norm of the average vector is zero")
	}

	// Normalize the average vector
	for i := range average {
		average[i] /= norm
	}

	return average, nil
}
