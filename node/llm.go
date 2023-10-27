package node

import (
	"fmt"

	"git.sr.ht/~primalmotion/simplai/llm"
	prompt "git.sr.ht/~primalmotion/simplai/prompt"
)

type LLMNode struct {
	llm llm.LLM
	*BaseNode
	options []llm.InferenceOption
}

func NewLLM(llm llm.LLM, options ...llm.InferenceOption) *LLMNode {
	return &LLMNode{
		BaseNode: New(),
		llm:      llm,
		options:  options,
	}
}

func (n *LLMNode) Execute(input prompt.Input) (string, error) {

	opts := n.options
	if stop := input.StopWords(); len(stop) > 0 {
		opts = append(opts, llm.OptionInferStop(stop...))
	}

	output, err := n.llm.Infer(input.Input(), opts...)
	if err != nil {
		return "", fmt.Errorf("unable to run llm inference: %w", err)
	}

	return n.BaseNode.Execute(prompt.NewInput(output))
}
