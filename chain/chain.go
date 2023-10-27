package chain

import (
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
)

type Chain struct {
	nodes []node.Node
}

func New(nodes ...node.Node) *Chain {

	chain := &Chain{
		nodes: nodes,
	}

	for i, n := range nodes {
		if len(nodes) > i+1 {
			n.Chain(nodes[i+1])
		}
	}

	return chain
}

func (c *Chain) Execute(input prompt.Input) (string, error) {

	return c.nodes[0].Execute(input)

}
