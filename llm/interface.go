package llm

// LLMN is the main interface to interfact with a LLM.
type LLM interface {
	Infer(text string, options ...InferenceOption) (string, error)
}
