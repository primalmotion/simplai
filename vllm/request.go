package vllm

import "fmt"

// VLLMRequest is the data holding the information to make a
// request to VLLM
type VLLMRequest struct {
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	Model            string         `json:"model,omitempty"`
	Prompt           string         `json:"prompt,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty"`
}

func (r VLLMRequest) String() string {
	return fmt.Sprintf(`---
model: %s
stop: %s
max_tokens: %d
temp: %f
top_p: %f
frequency_penalty: %f
presence_penalty: %f
logprobs: %d
logit_bias: %v

%s
---
`,
		r.Model,
		r.Stop,
		r.MaxTokens,
		r.Temperature,
		r.TopP,
		r.FrequencyPenalty,
		r.PresencePenalty,
		r.LogProbs,
		r.LogitBias,
		r.Prompt,
	)
}
