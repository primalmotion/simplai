package node

import (
	"context"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type Info struct {
	Name        string
	Description string
	Parameters  string
}

func LogNode(n Node, color string, format string, kwargs ...any) { // lulz
	render.Box(
		fmt.Sprintf("[%s]\n\n", n.Info().Name)+fmt.Sprintf(format, kwargs...),
		color,
	)
}

type Node interface {
	Info() Info
	Chain(Node)
	Next() Node
	Execute(context.Context, Input) (string, error)
}

type BaseNode struct {
	next Node
	desc Info
}

func New(desc Info) *BaseNode {
	return &BaseNode{
		desc: desc,
	}
}

func (n *BaseNode) Info() Info {
	return n.desc
}

func (n *BaseNode) Chain(next Node) {
	if n.next != nil {
		panic(fmt.Sprintf("node %s is already chained to %s", n.Info().Name, n.next.Info().Name))
	}
	n.next = next
}

func (n *BaseNode) Next() Node {
	return n.next
}

func (n *BaseNode) Execute(ctx context.Context, input Input) (string, error) {

	var err error
	var output string

	next := n.Next()
	if next != nil {
		output, err = next.Execute(ctx, input)
	} else {
		output = input.Input()
	}

	return output, err
}
