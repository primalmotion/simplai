package llm

import (
	tiktoken "github.com/pkoukk/tiktoken-go"
)

// CountToken is computing the max_tokens
// based on the model and the input.
func CountTokens(model, text string) int {
	e, err := tiktoken.EncodingForModel(model)
	if err != nil {
		e, err = tiktoken.GetEncoding("gpt2")
		if err != nil {
			return len([]rune(text))
		}
	}
	return len(e.Encode(text, nil, nil))
}
