package node

import (
	"git.sr.ht/~primalmotion/simplai/prompt"
)

type Node interface {
	Chain(next Node) Node
	Next() Node
	Execute(input prompt.Input) (string, error)
}

type BaseNode struct {
	next Node
}

func New() *BaseNode {
	return &BaseNode{}
}

func (n *BaseNode) Chain(next Node) Node {
	n.next = next
	return next
}

func (n *BaseNode) Next() Node {
	return n.next
}

func (n *BaseNode) Execute(input prompt.Input) (string, error) {

	next := n.Next()
	if next == nil {
		return input.Input(), nil
	}

	return next.Execute(input)
}
