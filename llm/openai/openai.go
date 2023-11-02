package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type openAIAPI struct {
	client  *http.Client
	url     *url.URL
	options options
	model   string
}

func New(api string, model string, opts ...Option) (*openAIAPI, error) {

	url, err := url.Parse(api)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url '%s': %w", api, err)
	}

	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	client := &http.Client{}
	return &openAIAPI{
		url:     url,
		model:   model,
		client:  client,
		options: o,
	}, nil
}

// Infer implements the node.Node interface
func (v *openAIAPI) Infer(ctx context.Context, prompt string, options ...llm.Option) (string, error) {

	config := v.options.defaultInferenceConfig
	config.Model = v.model
	config.MaxTokens = llm.CountTokens(v.model, prompt)

	for _, opt := range options {
		opt(&config)
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)

	vllmreq := request{
		LogitBias:         config.LogitBias,
		Model:             config.Model,
		Prompt:            prompt,
		Stop:              config.Stop,
		MaxTokens:         config.MaxTokens,
		Temperature:       config.Temperature,
		TopP:              config.TopP,
		TopK:              config.TopK,
		FrequencyPenalty:  config.FrequencyPenalty,
		RepetitionPenalty: config.RepetitionPenalty,
		PresencePenalty:   config.PresencePenalty,
		LogProbs:          config.LogProbs,
	}

	if config.Debug {
		render.Box(fmt.Sprintf("[openai-engine-request]\n\n%s", vllmreq), "4")
	}

	if err := encoder.Encode(vllmreq); err != nil {
		return "", fmt.Errorf("unable to encode request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/completions", v.url), buffer)
	if err != nil {
		return "", fmt.Errorf("unable to prepare request: %w", err)
	}

	resp, err := v.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("unable to send request: %w", err)
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server was unable to process the request: %s\n\n%s", resp.Status, content)
	}

	vllmresp := &response{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(vllmresp); err != nil {
		return "", fmt.Errorf("unable to decode the response: %w", err)
	}

	if config.Debug {
		render.Box(fmt.Sprintf("[openai-engine-response]\n\n%s", vllmresp), "4")
	}

	output := vllmresp.Choices[0].Text
	output = strings.TrimSpace(output)

	return output, nil
}
