package token

import tiktoken "github.com/pkoukk/tiktoken-go"

// Count is computing the max_tokens
// based on the model and the input.
func Count(model, text string) int {
	e, err := tiktoken.EncodingForModel(model)
	if err != nil {
		e, err = tiktoken.GetEncoding("gpt2")
		if err != nil {
			return len([]rune(text))
		}
	}
	return len(e.Encode(text, nil, nil))
}
