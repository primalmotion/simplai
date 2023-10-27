package chain

import (
	"git.sr.ht/~primalmotion/simplai/node"
)

func New(nodes ...node.Node) node.Node {

	for i, n := range nodes {

		if len(nodes) > i+1 {
			n.Chain(nodes[i+1])
		}
	}

	return nodes[0]
}
