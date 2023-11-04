package engine

import "context"

// LLM is the main interface to interact with a LLM.
type LLM interface {
	Infer(ctx context.Context, prompt string, options ...LLMOption) (string, error)
}
