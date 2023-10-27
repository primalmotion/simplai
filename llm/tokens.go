package llm

import (
	tiktoken "github.com/pkoukk/tiktoken-go"
)

func CountTokens(model, text string) int {
	e, err := tiktoken.EncodingForModel(model)
	if err != nil {
		e, err = tiktoken.GetEncoding("gpt2")
		if err != nil {
			return len([]rune(text)) / 4
		}
	}
	return len(e.Encode(text, nil, nil))
}
