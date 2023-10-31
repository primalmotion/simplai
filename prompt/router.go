package prompt

import (
	"context"
	"encoding/json"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/node"
)

const routerTemplate = `{{ .Input }}`

type routerInstruction struct {
	Action string `json:"action"`
	Params string `json:"params,omitempty"`
}

type Router struct {
	*node.Prompt
	subchainMap map[string]node.Node
}

func NewRouter(subchains ...node.Node) *Router {

	subchainMap := map[string]node.Node{}
	for _, s := range subchains {
		subchainMap[s.Name()] = s
	}
	return &Router{
		subchainMap: subchainMap,
		Prompt: node.NewPrompt(routerTemplate).
			WithName("router").
			WithDescription("route the input to a particular chain.").(*node.Prompt),
	}
}

func (n *Router) WithName(name string) node.Node {
	n.Prompt.WithName(name)
	return n
}

func (n *Router) WithDescription(desc string) node.Node {
	n.Prompt.WithDescription(desc)
	return n
}

func (n *Router) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Router) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (n *Router) isValidAction(action string) bool {

	if action == "" {
		return true
	}

	for k := range n.subchainMap {
		if action == k {
			return true
		}
	}
	return false
}

func (n *Router) Execute(ctx context.Context, in node.Input) (string, error) {

	inst := routerInstruction{}
	if err := json.Unmarshal([]byte(in.Input()), &inst); err != nil {
		return "", err
	}

	if !n.isValidAction(inst.Action) {
		return "", fmt.Errorf("invalid action name %s", inst.Action)
	}

	subchain := n.subchainMap[inst.Action]
	output, err := subchain.Execute(ctx, node.NewInput(inst.Params))
	if err != nil {
		return "", fmt.Errorf("unable to run subchain: %w", err)
	}

	return n.Prompt.Execute(ctx, node.NewInput(output))
}
