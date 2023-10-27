package vllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm"
)

type VLLM struct {
	client      *http.Client
	url         string
	model       string
	temperature float64
}

func NewVLLM(url string, model string, temperature float64) *VLLM {
	client := &http.Client{}
	return &VLLM{
		url:         url,
		model:       model,
		temperature: temperature,
		client:      client,
	}
}

func (v *VLLM) Infer(prompt string, options ...llm.InferenceOption) (string, error) {

	config := llm.NewInferenceConfig()
	config.Model = v.model
	config.Temperature = v.temperature

	for _, opt := range options {
		opt(&config)
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)

	vllmreq := VLLMRequest{
		LogitBias:        config.LogitBias,
		Model:            config.Model,
		Prompt:           prompt,
		Stop:             config.Stop,
		MaxTokens:        llm.CountTokens(v.model, prompt),
		Temperature:      config.Temperature,
		TopP:             config.TopP,
		FrequencyPenalty: config.FrequencyPenalty,
		PresencePenalty:  config.PresencePenalty,
		LogProbs:         config.LogProbs,
	}

	if err := encoder.Encode(vllmreq); err != nil {
		return "", fmt.Errorf("unable to encode request: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/completions", v.url), buffer)
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

	vllmresp := &VLLMResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(vllmresp); err != nil {
		return "", fmt.Errorf("Unable to decode the response: %w", err)
	}

	output := vllmresp.Choices[0].Text
	output = strings.TrimSpace(output)

	return output, nil
}
