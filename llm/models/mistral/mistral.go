package mistral

import (
	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

func NewLLM(llm llm.LLM, options ...llm.Option) *node.LLM {
	return node.NewLLM(
		node.Desc{
			Name: "mistral-llm",
		},
		llm,
		options...,
	)
}

func NewChatMemory() *node.ChatMemory {
	return node.NewChatMemory(
		node.Desc{
			Name: "mistral-memory",
		},
		"<|system|>",
		"<|assistant|>",
		"<|user|>",
	)
}
