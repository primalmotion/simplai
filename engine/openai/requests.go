package openai

import (
	"bytes"
	"encoding/json"
)

// request is the data holding the information to make a
// request to VLLM
type request struct {
	LogitBias         map[string]int `json:"logit_bias,omitempty"`
	Model             string         `json:"model,omitempty"`
	Prompt            string         `json:"prompt,omitempty"`
	Stop              []string       `json:"stop,omitempty"`
	MaxTokens         int            `json:"max_tokens,omitempty"`
	Temperature       float64        `json:"temperature"`
	TopP              float64        `json:"top_p"`
	TopK              int            `json:"top_k"`
	FrequencyPenalty  float64        `json:"frequency_penalty"`
	PresencePenalty   float64        `json:"presence_penalty"`
	RepetitionPenalty float64        `json:"repeat_penalty"`
	LogProbs          int            `json:"logprobs,omitempty"`
}

// String representation of request
func (r request) String() string {
	r.Prompt = ""

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(r)
	return buf.String()
}

// embeddingRequest is the openai structure to make
// an embedding request.
type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// String representation of embeddingRequest
func (r embeddingRequest) String() string {
	r.Input = []string{}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(r)
	return buf.String()
}

// rerankingRequest is the custom rerankingRequest
type rerankingRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

// String representation of reranking request.
func (r rerankingRequest) String() string {

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(r)
	return buf.String()
}
