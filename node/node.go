package node

import (
	"context"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/utils/render"
)

type Desc struct {
	Name        string
	Description string
}

func LogNode(n Node, color string, format string, kwargs ...any) { // lulz
	render.Box(
		fmt.Sprintf("[%s]\n\n", n.Desc().Name)+fmt.Sprintf(format, kwargs...),
		color,
	)
}

type PreHook func(Node, Input) (Input, error)
type PostHook func(Node, string) (string, error)

type Node interface {
	Desc() Desc
	Chain(Node) Node
	Next() Node
	WithPreHook(PreHook) Node
	WithPostHook(PostHook) Node
	Execute(context.Context, Input) (string, error)
}

type BaseNode struct {
	next     Node
	preHook  PreHook
	postHook PostHook
	desc     Desc
}

func New(desc Desc) *BaseNode {
	return &BaseNode{
		desc: desc,
	}
}

func (n *BaseNode) Desc() Desc {
	return n.desc
}

func (n *BaseNode) Chain(next Node) Node {
	if n.next != nil {
		panic(fmt.Sprintf("node %s is already chained to %s", n.Desc().Name, n.next.Desc().Name))
	}
	n.next = next
	return next
}

func (n *BaseNode) Next() Node {
	return n.next
}

func (n *BaseNode) Execute(ctx context.Context, input Input) (string, error) {

	var err error
	var output string

	if n.preHook != nil {
		input, err = n.preHook(n, input)
		if err != nil {
			return "", fmt.Errorf("error during pre hook: %w", err)
		}
	}

	next := n.Next()
	if next != nil {
		output, err = next.Execute(ctx, input)
	} else {
		output = input.Input()
	}

	if n.postHook != nil {
		output, err = n.postHook(n, output)
		if err != nil {
			return "", fmt.Errorf("error during post hook: %w", err)
		}
	}

	return output, err
}

func (n *BaseNode) WithPreHook(hook PreHook) Node {
	n.preHook = hook
	return n
}

func (n *BaseNode) WithPostHook(hook PostHook) Node {
	n.postHook = hook
	return n
}
