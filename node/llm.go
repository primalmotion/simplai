package node

import (
	"context"
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
		BaseNode: New().WithName("llm").(*BaseNode),
		llm:      llm,
		options:  options,
	}
}

func (n *LLM) WithName(name string) Node {
	n.BaseNode.WithName(name)
	return n
}

func (n *LLM) WithDescription(desc string) Node {
	n.BaseNode.WithDescription(desc)
	return n
}

func (n *LLM) WithPreHook(h PreHook) Node {
	n.BaseNode.WithPreHook(h)
	return n
}

func (n *LLM) WithPostHook(h PostHook) Node {
	n.BaseNode.WithPostHook(h)
	return n
}

func (n *LLM) Execute(ctx context.Context, input Input) (string, error) {

	output, err := n.llm.Infer(ctx,
		input.Input(),
		append(n.options, input.Options()...)...,
	)
	if err != nil {
		return "", fmt.Errorf("unable to run llm inference: %w", err)
	}

	return n.BaseNode.Execute(ctx, input.Derive(output))
}
