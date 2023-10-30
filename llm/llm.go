package llm

import "context"

// LLMN is the main interface to interfact with a LLM.
type LLM interface {
	Infer(ctx context.Context, prompt string, options ...Option) (string, error)
}
