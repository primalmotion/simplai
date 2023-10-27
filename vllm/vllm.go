package vllm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

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
		Prompt:      prompt,
		Model:       config.Model,
		MaxTokens:   int(math.Max(float64(llm.CountTokens(v.model, prompt)), 2048.0)),
		Temperature: config.Temperature,
	}

	fmt.Println(vllmreq)

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
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server was unable to process the request: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	vllmresp := &VLLMResponse{}
	if err := dec.Decode(vllmresp); err != nil {
		return "", fmt.Errorf("Unable to decode the response: %w", err)
	}

	return vllmresp.Choices[0].Text, nil
}
