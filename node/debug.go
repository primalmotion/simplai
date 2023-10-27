package node

import (
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type DebugNode struct {
	*BaseNode
}

func NewDebug() *DebugNode {
	return &DebugNode{
		BaseNode: New(),
	}
}

func (n *DebugNode) Execute(input prompt.Input) (string, error) {

	render.Box(input.Input(), "3")
	return n.BaseNode.Execute(input)
}
