package node

import (
	"fmt"

	"git.sr.ht/~primalmotion/simplai/llm"
)

type LLM struct {
	llm llm.LLM
	*BaseNode
	options []llm.Option
}

func NewLLM(llm llm.LLM, options ...llm.Option) *LLM {
	return &LLM{
		BaseNode: New(),
		llm:      llm,
		options:  options,
	}
}

func (n *LLM) Name() string {
	return "llm"
}

func (n *LLM) WithPreHook(h PreHook) Node {
	n.BaseNode.WithPreHook(h)
	return n
}

func (n *LLM) WithPostHook(h PostHook) Node {
	n.BaseNode.WithPostHook(h)
	return n
}

func (n *LLM) Execute(input Input) (string, error) {

	output, err := n.llm.Infer(
		input.Input(),
		append(n.options, input.Options()...)...,
	)
	if err != nil {
		return "", fmt.Errorf("unable to run llm inference: %w", err)
	}

	return n.BaseNode.Execute(
		NewInput(
			output,
			input.Options()...,
		),
	)
}
