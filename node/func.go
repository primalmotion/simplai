package node

import (
	"context"
)

// A Func is a node that will run a given function as its
// execution method.
type Func struct {
	executor func(context.Context, Input, Node) (string, error)
	*BaseNode
}

// NewFunc returns a new Func node that will use the provided function
// during it's execution.
func NewFunc(info Info, executor func(context.Context, Input, Node) (string, error)) *Func {
	return &Func{
		executor: executor,
		BaseNode: New(info),
	}
}

// Execute implements the Node interface.
func (n *Func) Execute(ctx context.Context, input Input) (string, error) {
	out, err := n.executor(ctx, input, n)
	if err != nil {
		return "", NewError(n, "unable to call executor: %w", err)
	}

	return n.BaseNode.Execute(ctx, input.WithInput(out))
}
