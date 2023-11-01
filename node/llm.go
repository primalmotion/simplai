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

func NewLLM(info Info, llm llm.LLM, options ...llm.Option) *LLM {
	return &LLM{
		BaseNode: New(info),
		llm:      llm,
		options:  options,
	}
}

func (n *LLM) Execute(ctx context.Context, input Input) (string, error) {

	opts := append(n.options, input.Options()...)
	if input.Debug() {
		opts = append(opts, llm.OptionDebug(true))
	}
	output, err := n.llm.Infer(ctx, input.Input(), opts...)
	if err != nil {
		return "", fmt.Errorf("unable to run llm inference: %w", err)
	}

	return n.BaseNode.Execute(ctx, input.Derive(output))
}
