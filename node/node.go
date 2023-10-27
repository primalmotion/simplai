package node

import (
	"fmt"

	"git.sr.ht/~primalmotion/simplai/llm"
	prompt "git.sr.ht/~primalmotion/simplai/prompt"
)

type Node struct {
	next    *Node
	prompt  prompt.Formatter
	llm     llm.LLM
	options []llm.InferenceOption
}

func New(
	llm llm.LLM,
	prmp prompt.Formatter,
	options ...llm.InferenceOption,
) *Node {
	return &Node{
		prompt:  prmp,
		llm:     llm,
		options: options,
	}
}

func (n *Node) Chain(next *Node) *Node {
	n.next = next
	return next
}

func (n *Node) Next() *Node {
	return n.next
}

func (n *Node) Execute(input prompt.Input) (string, error) {

	fprompt, err := n.prompt.Format(input)
	if err != nil {
		return "", fmt.Errorf("unable to format prompt: %w", err)
	}

	opts := append(n.options, llm.OptionInferStop(n.prompt.StopWords()...))

	output, err := n.llm.Infer(fprompt, opts...)
	if err != nil {
		return "", fmt.Errorf("unable to run llm inference: %w", err)
	}

	next := n.Next()
	if next == nil {
		return output, nil
	}

	return next.Execute(prompt.NewInput(output))
}
