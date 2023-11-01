package openai

import (
	"bytes"
	"encoding/json"
)

// request is the data holding the information to make a
// request to VLLM
type request struct {
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

func (r request) String() string {
	r.Prompt = ""

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(r)
	return buf.String()
}
