package node

import (
	"context"
	"fmt"
)

type Func struct {
	executor func(context.Context, Input, Node) (string, error)
	*BaseNode
}

func NewFunc(desc Desc, executor func(context.Context, Input, Node) (string, error)) *Func {
	return &Func{
		executor: executor,
		BaseNode: New(desc),
	}
}

func (n *Func) Execute(ctx context.Context, input Input) (string, error) {
	out, err := n.executor(ctx, input, n)
	if err != nil {
		return "", fmt.Errorf("[%s] unable to call executor func: %w", n.Desc().Name, err)
	}

	return n.BaseNode.Execute(ctx, input.Derive(out))
}