package vllm

import "fmt"

// VLLMRequest is the data holding the information to make a
// request to VLLM
type VLLMRequest struct {
	Model       string  `json:"model,omitempty"`
	Prompt      string  `json:"prompt,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

func (r VLLMRequest) String() string {
	return fmt.Sprintf(
		"-----\nmodel: %s\nmax: %d\ntemp: %f\n\n%s\n-----\n",
		r.Model,
		r.MaxTokens,
		r.Temperature,
		r.Prompt,
	)
}
