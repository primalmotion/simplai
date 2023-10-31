package prompt

import (
	"context"
	"encoding/json"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/node"
)

const routerTemplate = `{{.Input}}`

var RouterInfo = node.Info{
	Name:        "router",
	Description: "use to route the input to a particular chain.",
}

type routerInstruction struct {
	Action string `json:"action"`
	Params string `json:"params,omitempty"`
}

type Router struct {
	*node.Prompt
	subchainMap map[string]node.Node
	defaultNode node.Node
}

func NewRouter(defaultNode node.Node, subchains ...node.Node) *Router {

	subchainMap := map[string]node.Node{}
	for _, s := range subchains {
		subchainMap[s.Info().Name] = s
	}

	return &Router{
		subchainMap: subchainMap,
		defaultNode: defaultNode,
		Prompt: node.NewPrompt(
			RouterInfo,
			routerTemplate,
		),
	}
}

func (n *Router) Chain(next node.Node) {
	for _, s := range n.subchainMap {
		s.Chain(next)
	}
	n.defaultNode.Chain(next)
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
			n.Info().Name,
			inst.Action,
		)
	}

	var subchain node.Node
	if inst.Action != "" {
		subchain = n.subchainMap[inst.Action]
	} else {
		subchain = n.defaultNode
	}

	if in.Debug() {
		node.LogNode(
			n,
			"13",
			"received: %s\nexecuting subchain: %s",
			in.Input(),
			subchain.Info().Name,
		)
	}

	output, err := subchain.Execute(ctx, in.Derive(inst.Params).WithOptions(n.Options()...))
	if err != nil {
		return "", fmt.Errorf(
			"[%s] unable to run subchain '%s': %w",
			n.Info().Name,
			subchain.Info().Name,
			err,
		)
	}

	return n.Prompt.Execute(ctx, in.Derive(output))
}
