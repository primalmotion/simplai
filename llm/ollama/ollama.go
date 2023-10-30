package ollama

import (
	"context"
	"errors"

	"git.sr.ht/~primalmotion/simplai/llm"
	ollamaclient "git.sr.ht/~primalmotion/simplai/llm/ollama/internal"
)

var (
	ErrEmptyResponse       = errors.New("no response")
	ErrIncompleteEmbedding = errors.New("no all input got emmbedded")
)

// LLM is a ollama LLM implementation.
type LLM struct {
	client  *ollamaclient.Client
	options options
}

// New creates a new ollama LLM implementation.
func New(opts ...Option) (*LLM, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	client, err := ollamaclient.NewClient(o.ollamaServerURL)
	if err != nil {
		return nil, err
	}

	return &LLM{client: client, options: o}, nil
}

// Generate implemente the generate interface for LLM.
func (o *LLM) Infer(ctx context.Context, prompt string, options ...llm.Option) (string, error) {

	opts := llm.InferenceConfig{}
	for _, opt := range options {
		opt(&opts)
	}

	// Load back CallOptions as ollamaOptions
	ollamaOptions := o.options.ollamaOptions
	ollamaOptions.NumPredict = opts.MaxTokens
	ollamaOptions.Temperature = float32(opts.Temperature)
	ollamaOptions.Stop = opts.Stop
	ollamaOptions.TopK = opts.TopK
	ollamaOptions.TopP = float32(opts.TopP)
	ollamaOptions.Seed = opts.Seed
	ollamaOptions.RepeatPenalty = float32(opts.RepetitionPenalty)
	ollamaOptions.FrequencyPenalty = float32(opts.FrequencyPenalty)
	ollamaOptions.PresencePenalty = float32(opts.PresencePenalty)

	// Override LLM model if set as llm.InferenceOption
	model := o.options.model
	if opts.Model != "" {
		model = opts.Model
	}

	req := &ollamaclient.GenerateRequest{
		Model:    model,
		System:   o.options.system,
		Prompt:   prompt,
		Template: o.options.customModelTemplate,
		Options:  ollamaOptions,
	}

	resp, err := o.client.Infer(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Response, nil
}
