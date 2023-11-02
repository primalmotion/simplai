package node

import (
	"context"
	"fmt"

	"github.com/primalmotion/simplai/utils/render"
)

// A Subchain is a Node that holds a separate set of Nodes.
//
// It can be considered as a single Node.
//
//	[node] -> [[node]->[node]->[node]] -> [node]
//
// It can be useful to handle a separate set of Input or llm.Options. Subchains
// can also be used in Router nodes, that will execute a certain Subchain based
// on a condition.
//
//	[classify] -> [llm] -> [router] ?-> [summarize] ->[llm1] ?-> [search] ->
//	[llm2] ?-> [generate] -> [llm3]
//
// Subchain embeds the BaseNode and can be used as any other node.
type Subchain struct {
	*BaseNode
	nodes []Node
}

// NewSubchain creates a new Subchain with the given Info and given Nodes. The
// function will chain all of them in order. The given Node must not be already
// chained or it will panic.
func NewSubchain(info Info, nodes ...Node) *Subchain {

	for i, n := range nodes {
		if len(nodes) > i+1 && nodes[i+1] != nil {
			n.Chain(nodes[i+1])
		}
	}

	return &Subchain{
		BaseNode: New(info),
		nodes:    nodes,
	}
}

// NewSubchainWithName is a convenience function to create a named Subchain.
func NewSubchainWithName(name string, nodes ...Node) *Subchain {
	return NewSubchain(Info{Name: name}, nodes...)
}

// Execute executes the Subchain. It will first execute all the internal
// nodes, then execute the receiver's Next() Node.
func (n *Subchain) Execute(ctx context.Context, input Input) (string, error) {

	if input.Debug() {
		render.Box(fmt.Sprintf("[%s]", n.info.Name), "3")
	}

	first := n.nodes[0]
	output, err := first.Execute(ctx, input)
	if err != nil {
		return "", NewError(n, "unable to execute node '%s': %w", first.Info().Name, err)
	}

	return n.BaseNode.Execute(ctx, input.WithInput(output))
}
