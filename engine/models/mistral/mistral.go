package mistral

import (
	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/node"
)

// NewLLM returns a new node.LLM named mistral-llm
func NewLLM(eng engine.LLM, options ...engine.LLMOption) *node.LLM {
	return node.NewLLM(
		node.Info{
			Name: "mistral-llm",
		},
		eng,
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
