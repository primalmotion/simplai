package node

import (
	"fmt"

	"git.sr.ht/~primalmotion/simplai/llm"
	prompt "git.sr.ht/~primalmotion/simplai/prompt"
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

func (n *LLM) Execute(input prompt.Input) (string, error) {

	opts := n.options
	if iopts := input.Options(); len(iopts) > 0 {
		opts = append(opts, iopts...)
	}

	output, err := n.llm.Infer(input.Input(), opts...)
	if err != nil {
		return "", fmt.Errorf("unable to run llm inference: %w", err)
	}

	return n.BaseNode.Execute(prompt.NewInput(output))
}
