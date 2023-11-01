package node

import (
	"context"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/utils/render"
)

// A Chain is a Node that holds a separate set of Nodes.
//
// It can be seen as a nested chain inside a chain, that can be considered as a
// single node.
//
//	[node] -> [[node]->[node]->[node]] -> [node]
//
// It can be useful to handle a separate set of Input llm.Options for instance.
// Chains can also be used in Router nodes, that will decide which subchain it
// should forward the Input to.
//
//	[classify] -> [llm] -> [router] ?-> [summarize] ->[llm1]
//									?-> [search] -> [llm2]
//									?-> [generate] -> [llm3]
//
// Chain embeds the BaseNode and can be used as any other node.
type Chain struct {
	*BaseNode
	nodes []Node
}

// NewChain creates a new chain with the given Info and given Nodes. The
// function will chain all of them in order. The given node must not be already
// chained or it will panic.
func NewChain(desc Info, nodes ...Node) *Chain {

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

// NewChainWithName is a convenience function to create a chain
// with a name as a string.
func NewChainWithName(name string, nodes ...Node) *Chain {
	return NewChain(Info{Name: name}, nodes...)
}

// Execute executes the chain. It will first execute all the internal
// nodes, then execute the receiver's Next() Node.
func (c *Chain) Execute(ctx context.Context, input Input) (string, error) {

	if input.Debug() {
		render.Box(fmt.Sprintf("[%s]", c.info.Name), "3")
	}

	output, err := c.nodes[0].Execute(ctx, input)
	if err != nil {
		return "", fmt.Errorf(
			"[%s] unable to execute node '%s': %w",
			c.info.Name,
			c.nodes[0].Info().Name,
			err,
		)
	}

	return c.BaseNode.Execute(ctx, input.Derive(output))
}
