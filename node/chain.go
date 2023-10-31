package node

import (
	"context"
	"fmt"
)

type Chain struct {
	*BaseNode
	nodes []Node
}

func NewChain(nodes ...Node) *Chain {

	for i, n := range nodes {
		if len(nodes) > i+1 && nodes[i+1] != nil {
			n.Chain(nodes[i+1])
		}
	}

	return &Chain{
		BaseNode: New().WithName("chain").(*BaseNode),
		nodes:    nodes,
	}
}

func (n *Chain) WithName(name string) Node {
	n.BaseNode.WithName(name)
	return n
}

func (n *Chain) WithDescription(desc string) Node {
	n.BaseNode.WithDescription(desc)
	return n
}

func (n *Chain) WithPreHook(h PreHook) Node {
	n.BaseNode.WithPreHook(h)
	return n
}

func (n *Chain) WithPostHook(h PostHook) Node {
	n.BaseNode.WithPostHook(h)
	return n
}

func (c *Chain) Execute(ctx context.Context, input Input) (string, error) {

	output, err := c.nodes[0].Execute(ctx, input)
	if err != nil {
		return "", fmt.Errorf("unable to execute chained node: %w", err)
	}

	return c.BaseNode.Execute(ctx, NewInput(output))
}
