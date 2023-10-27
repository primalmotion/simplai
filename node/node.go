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

func New(llm llm.LLM, prmp prompt.Formatter, options ...llm.InferenceOption) *Node {
	return &Node{
		prompt:  prmp,
		llm:     llm,
		options: options,
	}
}

func (n *Node) AddNode(next *Node) *Node {
	n.next = next
	return next
}

func (n *Node) Execute(input prompt.Input) (string, error) {

	fprompt, err := n.prompt.Format(input)
	if err != nil {
		return "", fmt.Errorf("unable to format prompt: %w", err)
	}

	fmt.Println("------------")
	fmt.Println("NODE PROMPT:")
	fmt.Println(fprompt)

	output, err := n.llm.Infer(fprompt, n.options...)
	if err != nil {
		return "", err
	}

	if n.next == nil {
		return output, nil
	}

	fmt.Println("OUTPUT:", output)
	return n.next.Execute(prompt.NewInput(output))
}
