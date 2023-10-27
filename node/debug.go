package node

import (
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type Debug struct {
	*BaseNode
}

func NewDebug() *Debug {
	return &Debug{
		BaseNode: New(),
	}
}

func (n *Debug) Name() string {
	return "llm"
}

func (n *Debug) Execute(input prompt.Input) (string, error) {

	render.Box(input.Input(), "3")
	return n.BaseNode.Execute(input)
}
