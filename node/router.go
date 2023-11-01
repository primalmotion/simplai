package node

import (
	"context"
	"encoding/json"
	"fmt"
)

type RouterDeciderFunc func(context.Context, Input, Node, map[string]Node) (Node, Input, error)

type Router struct {
	defaultChain Node
	*BaseNode
	chains  map[string]Node
	decider RouterDeciderFunc
}

func NewRouter(
	info Info,
	decider RouterDeciderFunc,
	defaultChain Node,
	chains ...Node,
) *Router {

	chs := make(map[string]Node, len(chains))
	for _, c := range chains {
		chs[c.Info().Name] = c
	}

	return &Router{
		decider:      decider,
		defaultChain: defaultChain,
		chains:       chs,
		BaseNode:     New(info),
	}
}

// Chain is overriding the default Node behavior. A Router never uses its
// internal Next(). Instead, Chain() will chain all the Router's subchains
// (including the default) to the given Node.
// When the router is executed, it will hand off execution to the selected
// subchain.
func (n *Router) Chain(next Node) {

	for _, c := range n.chains {
		c.Chain(next)
	}

	n.defaultChain.Chain(next)
}

func (n *Router) Execute(ctx context.Context, input Input) (string, error) {

	selected, dinput, err := n.decider(ctx, input, n.defaultChain, n.chains)
	if err != nil {
		return "", NewError(n, "unable to run decider func: %w", err)
	}

	output, err := selected.Execute(ctx, dinput)
	if err != nil {
		return "", NewError(n, "unable to run chain '%s': %w", selected.Info().Name, err)
	}

	return n.BaseNode.Execute(ctx, input.Derive(output))
}

type RouterSimpleInput struct {
	Params map[string]any `json:"params,omitempty"`
	Name   string         `json:"name"`
	Input  string         `json:"input,omitempty"`
}

func RouterSimpleDeciderFunc(
	ctx context.Context,
	input Input,
	def Node,
	chains map[string]Node,
) (Node, Input, error) {

	rinput := RouterSimpleInput{}
	if err := json.Unmarshal([]byte(input.Input()), &rinput); err != nil {
		return nil, input, fmt.Errorf("unable to unmarshal input '%s': %w", input.Input(), err)
	}

	if selected := chains[rinput.Name]; selected != nil {
		return selected, input.Derive(rinput.Input), nil
	}

	return def, input.Derive(rinput.Input), nil
}
