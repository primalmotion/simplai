package llm

import "context"

// LLMN is the main interface to interact with a LLM.
type LLM interface {
	Infer(ctx context.Context, prompt string, options ...Option) (string, error)
}
