package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/engine/internal/embedding"
	"github.com/primalmotion/simplai/engine/internal/token"
	"github.com/primalmotion/simplai/utils/render"
)

// check interface compliance.
var (
	_ engine.Embedder = &openAIAPI{}
	_ engine.LLM      = &openAIAPI{}
)

type openAIAPI struct {
	client  *http.Client
	url     *url.URL
	model   string
	options options
}

// New Return a new openAIAPI client.
func New(api string, model string, opts ...Option) (*openAIAPI, error) { //nolint:revive

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
func (v *openAIAPI) Infer(ctx context.Context, prompt string, options ...engine.LLMOption) (string, error) {

	config := v.options.defaultInferenceConfig
	config.Model = v.model
	config.MaxTokens = token.Count(v.model, prompt)

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

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

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

// Embed implements the Embedder interface.
func (v *openAIAPI) Embed(ctx context.Context, chunks []string, options ...engine.EmbeddingOption) ([][]float64, error) {

	config := defaultEmbeddingConfig()
	for _, opt := range options {
		opt(&config)
	}

	model := config.Model
	if model == "" {
		model = v.model
	}

	emb := make([][]float64, 0, len(chunks))

	batches := embedding.Batch(chunks, config.BatchSize)
	for _, batch := range batches {

		currentEmbeddings := [][]float64{}

		buffer := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buffer)

		req := &embeddingRequest{
			Model: model,
			Input: batch,
		}

		if config.Debug {
			render.Box(fmt.Sprintf("[openai-embedding-request]\n\n%s", req), "4")
		}

		if err := encoder.Encode(req); err != nil {
			return nil, fmt.Errorf(
				"unable to encode request: %w", err)
		}

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/embeddings", v.url), buffer)
		if err != nil {
			return nil, fmt.Errorf("unable to prepare request: %w", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")

		resp, err := v.client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("unable to send request: %w", err)
		}

		defer func() {
			if resp.Body != nil {
				_ = resp.Body.Close()
			}
		}()

		if resp.StatusCode != http.StatusOK {
			content, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("server was unable to process the request: %s\n\n%s", resp.Status, content)
		}

		embResp := &embeddingResponse{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(embResp); err != nil {
			return nil, fmt.Errorf("unable to decode the response: %w", err)
		}

		if config.Debug {
			render.Box(fmt.Sprintf("[openai-embedding-response]\n\n%s", embResp), "4")
		}

		if len(embResp.Data) == 0 {
			return nil, errors.New("empty response")
		}

		for i := 0; i < len(embResp.Data); i++ {
			currentEmbeddings = append(currentEmbeddings, embResp.Data[i].Embedding)
		}

		if len(batch) != len(currentEmbeddings) {
			return currentEmbeddings, errors.New("no all input got emmbedded")
		}

		// get num of token in that batch
		// we should use the encoder of the model to get the tokens
		// but its not available. So we fall back on tiktoken
		numTokens := make([]float64, 0, len(batch))
		for _, text := range batch {
			numTokens = append(numTokens, float64(token.Count(config.Model, text)))
		}

		if len(currentEmbeddings) > 1 {
			combinedVectors, err := embedding.CombineBatchedEmbedding(currentEmbeddings, numTokens)
			if err != nil {
				return [][]float64{}, err
			}
			emb = append(emb, combinedVectors)
			continue
		}

		emb = append(emb, currentEmbeddings...)
	}

	return emb, nil
}
