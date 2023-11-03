package openai

import (
	"bytes"
	"encoding/json"
)

// VLLMResponse is the structure describing a VLLM response.

type response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model,omitempty"`
	Choices []struct {
		Text         string `json:"text,omitempty"`
		FinishReason string `json:"finish_reason,omitempty"`
		Index        int    `json:"index,omitempty"`
	} `json:"choices,omitempty"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens,omitempty"`
		TotalTokens      int `json:"total_tokens,omitempty"`
		CompletionTokens int `json:"completion_tokens,omitempty"`
	} `json:"usage,omitempty"`
}

// String representation of response.
func (r response) String() string {
	r.Choices = nil
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(r)
	return buf.String()
}

// embeddingResponse is the openai struct to get
// an embedding response.
type embeddingResponse struct {
	Object string `json:"object"`
	Model  string `json:"model"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// String representation of embeddingResponse.
func (r embeddingResponse) String() string {
	r.Data = nil
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(r)
	return buf.String()
}
