package node

import (
	"fmt"

	"git.sr.ht/~primalmotion/simplai/prompt"
)

type PreHook func(Node, prompt.Input) (prompt.Input, error)
type PostHook func(Node, string) (string, error)

type Node interface {
	Chain(next Node) Node
	Next() Node
	Execute(input prompt.Input) (string, error)
	Name() string
	WithPreHook(PreHook) Node
	WithPostHook(PostHook) Node
}

type BaseNode struct {
	next     Node
	preHook  PreHook
	postHook PostHook
}

func New() *BaseNode {
	return &BaseNode{}
}

func (n *BaseNode) Name() string {
	return "base"
}

func (n *BaseNode) Chain(next Node) Node {
	n.next = next
	return next
}

func (n *BaseNode) Next() Node {
	return n.next
}

func (n *BaseNode) Execute(input prompt.Input) (string, error) {

	var err error
	var output string

	if n.preHook != nil {
		input, err = n.preHook(n, input)
		if err != nil {
			return "", fmt.Errorf("error during pre hook: %w", err)
		}
	}

	next := n.Next()

	if next != nil {
		output, err = next.Execute(input)
	} else {
		output = input.Input()
	}

	if n.postHook != nil {
		output, err = n.postHook(n, output)
		if err != nil {
			return "", fmt.Errorf("error during post hook: %w", err)
		}
	}

	return output, err
}

func (n *BaseNode) WithPreHook(hook PreHook) Node {
	n.preHook = hook
	return n
}

func (n *BaseNode) WithPostHook(hook PostHook) Node {
	n.postHook = hook
	return n
}
