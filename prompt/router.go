package prompt

import (
	"context"
	"encoding/json"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/node"
)

const routerTemplate = `{{.Input}}`

var RouterDesc = node.Desc{
	Name:        "router",
	Description: "route the input to a particular chain.",
}

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
		subchainMap[s.Desc().Name] = s
	}

	return &Router{
		subchainMap: subchainMap,
		Prompt: node.NewPrompt(
			RouterDesc,
			routerTemplate,
		),
	}
}

func (n *Router) Chain(next node.Node) node.Node {
	for _, s := range n.subchainMap {
		s.Chain(next)
	}
	return next
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
		return "", fmt.Errorf(
			"[%s] invalid action name '%s'",
			n.Desc().Name,
			inst.Action,
		)
	}

	subchain := n.subchainMap[inst.Action]

	if in.Debug() {
		node.LogNode(
			n,
			"13",
			"received: %s\nexecuting subchain: %s",
			in.Input(),
			subchain.Desc().Name,
		)
	}

	output, err := subchain.Execute(
		ctx,
		in.
			Derive(inst.Params).
			WithOptions(n.Options()...),
	)
	if err != nil {
		return "", fmt.Errorf(
			"[%s] unable to run subchain '%s': %w",
			n.Desc().Name,
			subchain.Desc().Name,
			err,
		)
	}

	// return output, nil
	return n.Prompt.Execute(ctx, in.Derive(output))
}
