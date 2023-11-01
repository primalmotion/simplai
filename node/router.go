package node

import (
	"context"
	"encoding/json"
)

// RouterDeciderFunc is the type of the function used by the Router. It is
// responsible for making a decision on which chains the router should send the
// input to. It is given the Router's context, the Input as well as a map of
// available subchains, keyed by their name.
// It must return which Node to use and which Input to pass to it.
type RouterDeciderFunc func(context.Context, Input, Node, map[string]Node) (Node, Input, error)

// A Router is a node that can route its Input to
// one of several chains. The decision is made by a RouterDeciderFunc.
type Router struct {
	defaultChain Node
	*BaseNode
	chains  map[string]Node
	decider RouterDeciderFunc
}

// NewRouter returns a new Router.
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

// Execute implements the Node interface.
func (n *Router) Execute(ctx context.Context, input Input) (string, error) {

	selected, dinput, err := n.decider(ctx, input, n.defaultChain, n.chains)
	if err != nil {
		return "", NewError(n, "unable to run decider func: %w", err)
	}

	output, err := selected.Execute(ctx, dinput)
	if err != nil {
		return "", NewError(n, "unable to run chain '%s': %w", selected.Info().Name, err)
	}

	return n.BaseNode.Execute(ctx, input.WithInput(output))
}

// RouterSimpleInput is the data structure that will be used
// to decide how to router the traffic by the RouterSimpleDeciderFunc.
type RouterSimpleInput struct {
	Params map[string]any `json:"params,omitempty"`
	Name   string         `json:"name"`
	Input  string         `json:"input,omitempty"`
}

// RouterSimpleDeciderFunc is a decider that will use an LLM output
// formatted as JSON structure that will be decoded by RouterSimpleInput.
func RouterSimpleDeciderFunc(
	ctx context.Context,
	input Input,
	def Node,
	chains map[string]Node,
) (Node, Input, error) {

	rinput := RouterSimpleInput{}
	if err := json.Unmarshal([]byte(input.Input()), &rinput); err != nil {
		return nil, input, NewPromptError(
			"I MUST write a valid JSON structure",
			err,
		)
	}

	if selected := chains[rinput.Name]; selected != nil {
		return selected, input.WithInput(rinput.Input), nil
	}

	return def, input.WithInput(rinput.Input), nil
}
