package mistral

import (
	"github.com/primalmotion/simplai/llm"
	"github.com/primalmotion/simplai/node"
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

// NewChatMemory a new node.ChatMemory named mistral-memory
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
