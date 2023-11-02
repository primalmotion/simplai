package mistral

import (
	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

// NewLLM returns a new node.LLM named mistral-llm
func NewLLM(llm llm.LLM, options ...llm.Option) *node.LLM {
	return node.NewLLM(
		node.Info{
			Name: "mistral-llm",
		},
		llm,
		options...,
	)
}

// NewLLM returns a new node.ChatMemory named mistral-memory
// correctly configured to optmized the chats based on
// mistral training.
func NewChatMemory() *node.ChatMemory {
	return node.NewChatMemory(
		node.Info{
			Name: "mistral-memory",
		},
		"<|system|>",
		"<|assistant|>",
		"<|user|>",
	)
}
