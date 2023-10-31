package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type OpenAIAPI struct {
	client      *http.Client
	url         string
	model       string
	temperature float64
}

func NewOpenAIAPI(url string, model string, temperature float64) *OpenAIAPI {
	client := &http.Client{}
	return &OpenAIAPI{
		url:         url,
		model:       model,
		temperature: temperature,
		client:      client,
	}
}

func (v *OpenAIAPI) Infer(ctx context.Context, prompt string, options ...llm.Option) (string, error) {

	config := llm.NewInferenceConfig()
	config.Model = v.model
	config.Temperature = v.temperature
	config.MaxTokens = llm.CountTokens(v.model, prompt)

	for _, opt := range options {
		opt(&config)
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)

	vllmreq := request{
		LogitBias:        config.LogitBias,
		Model:            config.Model,
		Prompt:           prompt,
		Stop:             config.Stop,
		MaxTokens:        config.MaxTokens,
		Temperature:      config.Temperature,
		TopP:             config.TopP,
		FrequencyPenalty: config.FrequencyPenalty,
		PresencePenalty:  config.PresencePenalty,
		LogProbs:         config.LogProbs,
	}

	if config.Debug {
		render.Box(vllmreq.String(), "8")
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
			resp.Body.Close()
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
		render.Box(vllmresp.String(), "8")
	}

	output := vllmresp.Choices[0].Text
	output = strings.TrimSpace(output)

	return output, nil
}
