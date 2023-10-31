package prompt

import (
	"context"
	"fmt"

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

var ClassifierDesc = node.Desc{
	Name:        "classifier",
	Description: "used to classify the intent of the user.",
}

type Classifier struct {
	*node.Prompt
	subchainMap map[string]node.Desc
}

func NewClassifier(subchains ...node.Desc) *Classifier {

	subchainMap := map[string]node.Desc{}
	for _, s := range subchains {
		subchainMap[s.Name] = s
	}

	return &Classifier{
		subchainMap: subchainMap,
		Prompt: node.NewPrompt(
			ClassifierDesc,
			classifierTemplate,
			llm.OptionStop("\n"),
			llm.OptionMaxTokens(100),
		),
	}
}

func (n *Classifier) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Classifier) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (n *Classifier) subchainNames() []string {
	out := make([]string, len(n.subchainMap))
	i := 0
	for _, v := range n.subchainMap {
		out[i] = fmt.Sprintf(`{"action": "%s"}`, v.Name)
		i++
	}
	return out
}

func (n *Classifier) Execute(ctx context.Context, in node.Input) (output string, err error) {
	return n.Prompt.Execute(
		ctx,
		in.WithKeyValue("subchains", n.subchainMap),
	)
}
