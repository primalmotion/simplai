package prompt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

const classifierTemplate = `classify the user input and map it to a known action name.
your output will be used by a machine and it needs to be perfect.
You must conform to the protocol described below at all times.

## PROTOCOL

Actions are the various classes you must discriminate.
They are described as <name>:<intent>.

Given a user input, you must find which action is describing the best the
user's intent then what are the parameters for this action.

For example:

	INPUT: write a hello world program in python
	ACTION: {"action": "code", "params": "hello world in python"}

	INPUT: compose a song about bananas
	ACTION: {"action": "compose", "params": "bananas"}

If the input not explicitely map to any known actions, you must exactly write:

	ACTION: {"action": ""}

Your answer MUST be a valid JSON.

It is VERY IMPORTANT you understand that you MUST follow this protocol at all
costs, no matter what, and at all circumstances.

## KNOWN ACTIONS
{{ range $k, $v := .Get "subchains" }}
	{{ $k }}: {{ $v.Description -}}
{{ end }}

Remember: ACTION must only be one of: {{ range $k, $v := .Get "subchains" }}{{$k}}, {{end}}
{{ if (.Get "scratchpad") }}
## OBSERVATION

{{ .Get "scratchpad" }} {{ end }}

## PROCEED

INPUT: {{ .Input }}
ACTION:`

type Classifier struct {
	*node.Prompt
	subchainMap map[string]node.Node
}

func NewClassifier(subchains ...node.Node) *Classifier {

	subchainMap := map[string]node.Node{}
	for _, s := range subchains {
		subchainMap[s.Name()] = s
	}
	return &Classifier{
		subchainMap: subchainMap,
		Prompt: node.NewPrompt(
			classifierTemplate,
			llm.OptionStop("\n"),
			llm.OptionMaxTokens(100),
		).
			WithName("classifier").
			WithDescription("used to classify the intent of the user.").(*node.Prompt),
	}
}

func (n *Classifier) WithName(name string) node.Node {
	n.Prompt.WithName(name)
	return n
}

func (n *Classifier) WithDescription(desc string) node.Node {
	n.Prompt.WithDescription(desc)
	return n
}

func (n *Classifier) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Classifier) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (n *Classifier) isValidAction(action string) bool {

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

func (n *Classifier) subchainNames() []string {
	out := make([]string, len(n.subchainMap))
	i := 0
	for _, v := range n.subchainMap {
		out[i] = fmt.Sprintf(`{"action": "%s"}`, v.Name())
		i++
	}
	return out
}

func (n *Classifier) Execute(ctx context.Context, in node.Input) (output string, err error) {

	var i int
	for i = 0; i <= 3; i++ {

		in := in.WithKeyValue("subchains", n.subchainMap)

		output, err = n.Prompt.Execute(ctx, in)
		if err != nil {
			return "", err
		}

		out := map[string]any{}
		if err := json.Unmarshal([]byte(output), &out); err != nil {
			in = in.WithKeyValue("scratchpad", "I failed to generate a valid json. I must generate a valid json.")
			continue
		}

		if !n.isValidAction(out["action"].(string)) {
			in = in.WithKeyValue(
				"scratchpad",
				fmt.Sprintf(
					`{"action": "%s"} is invalid. Only use one of: %s. If nothing matches, then write {"action": ""}`,
					out["action"],
					strings.Join(n.subchainNames(), ", "),
				),
			)
			continue
		}

		break
	}

	if i >= 3 {
		return `{"action": ""}`, nil
	}

	return output, nil
}
