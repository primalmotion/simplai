package chain

import (
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
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
		if len(nodes) > i+1 {
			if nodes[i+1] != nil {
				n.Chain(nodes[i+1])
			}
		}
	}

	return &Chain{
		nodes: internalNodes,
	}
}

func (c *Chain) Execute(input prompt.Input) (string, error) {

	return c.nodes[0].Execute(input)

}
