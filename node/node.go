package node

import (
	"context"
	"fmt"
)

type PreHook func(Node, Input) (Input, error)
type PostHook func(Node, string) (string, error)

type Node interface {
	Chain(next Node) Node
	Next() Node

	WithName(string) Node
	Name() string
	WithDescription(string) Node
	Description() string
	WithPreHook(PreHook) Node
	WithPostHook(PostHook) Node

	Execute(ctx context.Context, input Input) (string, error)
}

type BaseNode struct {
	next        Node
	preHook     PreHook
	postHook    PostHook
	name        string
	description string
}

func New() *BaseNode {
	return &BaseNode{}
}

func (n *BaseNode) WithName(name string) Node {
	n.name = name
	return n
}

func (n *BaseNode) WithDescription(desc string) Node {
	n.description = desc
	return n
}

func (n *BaseNode) Name() string {
	return n.name
}

func (n *BaseNode) Description() string {
	return n.description
}

func (n *BaseNode) Chain(next Node) Node {
	n.next = next
	return next
}

func (n *BaseNode) Next() Node {
	return n.next
}

func (n *BaseNode) Execute(ctx context.Context, input Input) (string, error) {

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
		output, err = next.Execute(ctx, input)
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
