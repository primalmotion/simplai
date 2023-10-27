package node

import (
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type PrintNode struct {
	*BaseNode
}

func NewPrintNode() *PrintNode {
	return &PrintNode{
		BaseNode: New(),
	}
}

func (n *PrintNode) Execute(input prompt.Input) (string, error) {

	render.Box(input.Input(), "3")
	return n.BaseNode.Execute(input)
}
