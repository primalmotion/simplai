package node

import (
	"context"
)

type Func struct {
	executor func(context.Context, Input, Node) (string, error)
	*BaseNode
}

func NewFunc(info Info, executor func(context.Context, Input, Node) (string, error)) *Func {
	return &Func{
		executor: executor,
		BaseNode: New(info),
	}
}

func (n *Func) Execute(ctx context.Context, input Input) (string, error) {
	out, err := n.executor(ctx, input, n)
	if err != nil {
		return "", NewError(n, "unable to call executor: %w", err)
	}

	return n.BaseNode.Execute(ctx, input.Derive(out))
}
