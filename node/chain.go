package node

import (
	"context"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type Chain struct {
	*BaseNode
	nodes []Node
}

func NewChainWithName(name string, nodes ...Node) *Chain {
	return NewChain(Desc{Name: name}, nodes...)
}

func NewChain(desc Desc, nodes ...Node) *Chain {

	for i, n := range nodes {
		if len(nodes) > i+1 && nodes[i+1] != nil {
			n.Chain(nodes[i+1])
		}
	}

	return &Chain{
		BaseNode: New(desc),
		nodes:    nodes,
	}
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

	if input.Debug() {
		render.Box(fmt.Sprintf("[%s]", c.desc.Name), "3")
	}

	output, err := c.nodes[0].Execute(ctx, input)
	if err != nil {
		return "", fmt.Errorf(
			"[%s] unable to execute node '%s': %w",
			c.desc.Name,
			c.nodes[0].Desc().Name,
			err,
		)
	}

	return c.BaseNode.Execute(ctx, input.Derive(output))
}
