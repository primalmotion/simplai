package ollama

import (
	"context"
	"fmt"
	"net/url"

	"github.com/primalmotion/simplai/llm"
	ollamaclient "github.com/primalmotion/simplai/llm/ollama/internal"
	"github.com/primalmotion/simplai/utils/render"
)

// LLM is a ollama LLM implementation.
type ollamaAPI struct {
	client  *ollamaclient.Client
	model   string
	options options
}

// New creates a new ollama LLM implementation.
func New(api string, model string, opts ...Option) (*ollamaAPI, error) { //nolint:revive

	url, err := url.Parse(api)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url '%s': %w", api, err)
	}

	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	return &ollamaAPI{
		client:  ollamaclient.NewClient(url),
		model:   model,
		options: o,
	}, nil
}

// Generate implemente the generate interface for LLM.
func (o *ollamaAPI) Infer(ctx context.Context, prompt string, options ...llm.Option) (string, error) {

	opts := o.options.defaultInferenceConfig
	opts.Model = o.model
	opts.MaxTokens = llm.CountTokens(o.model, prompt)

	for _, opt := range options {
		opt(&opts)
	}

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

	req := &ollamaclient.GenerateRequest{
		Model:    opts.Model,
		System:   o.options.system,
		Prompt:   prompt,
		Template: o.options.customModelTemplate,
		Options:  ollamaOptions,
		Raw:      o.options.raw,
	}

	if opts.Debug {
		render.Box(fmt.Sprintf("[ollama-engine-request]\n\n%s", req), "4")
	}

	resp, err := o.client.Infer(ctx, req)
	if err != nil {
		return "", err
	}

	if opts.Debug {
		render.Box(fmt.Sprintf("[ollama-engine-response]\n\n%s", resp), "4")
	}

	return resp.Response, nil
}
