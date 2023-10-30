package chain

import (
	"context"

	"git.sr.ht/~primalmotion/simplai/node"
)

type Chain struct {
	nodes []node.Node
}

func New(nodes ...node.Node) *Chain {

	if len(nodes) == 0 {
		return &Chain{}
	}

	internalNodes := []node.Node{}

	for i, n := range nodes {
		if n == nil {
			continue
		}
		internalNodes = append(internalNodes, n)
		if len(nodes) > i+1 && nodes[i+1] != nil {
			n.Chain(nodes[i+1])
		}
	}

	return &Chain{
		nodes: internalNodes,
	}
}

func (c *Chain) Execute(ctx context.Context, input node.Input) (string, error) {

	return c.nodes[0].Execute(ctx, input)

}
