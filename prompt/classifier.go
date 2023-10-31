package prompt

import (
	"context"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

const classifierTemplate = `classify the user input and map it to a known tool name.
your output will be used by a machine and it needs to be perfect.
You must conform to the protocol described below at all times.

## PROTOCOL

You must identify which of the known tools is the best fullfill the user request.
They are described as:

	NAME: <name>
	USAGE: <intent>
	PARAMS: <information about parameters>

The tool NAME is what you must classify. You will use the USAGE to help you choose a tool.
You also need to pay close attention to the PARAMS, if any, in order to complete the paramaters.

You MUST write a valid JSON output, describing the tools to use in the form:

	{"action":"<tool-name>","parameters":"<required-tool-parameters"}

For example, given the tools:

	- NAME: code
	  USAGE: use to write some code in various programming language
	  PARAMS: the description of the desired code to produce

	- NAME: compose
	  USAGE: use to compose a song
	  PARAMS: the subject the song is about.

Here is some output example:

	INPUT: write a hello world program in python
	ACTION: {"action":"code","params":"hello world in python"}

	INPUT: compose a song about bananas
	ACTION: {"action":"compose","params":"bananas"}

	INPUT: jump over the bridge
	ACTION: {"action":"","params":"jump over the bridge"}

If the input not explicitely map to any known actions, you MUST exactly write:

	ACTION: {"action":"","params":"{{.Input}}"}

It is VERY IMPORTANT you remember that you MUST follow this protocol no matter
what, in all circumstances.

## AVAILABLE TOOLS
{{ range $k, $v := .Get "subchains" }}
	- NAME: {{ $k }}
	  USAGE: {{ $v.Description }}
	  PARAMS: {{ $v.Parameters }}
{{ end }}

Remember: ACTION must only be one of: {{ range $k, $v := .Get "subchains" }}{{$k}}, {{end}}
Pay attention to the action description if it details what the params should be.

{{ if (.Get "scratchpad") }}
## OBSERVATION

{{ .Get "scratchpad" }} {{ end }}

## PROCEED

INPUT: {{ .Input }}
ACTION:`

var ClassifierInfo = node.Info{
	Name:        "classifier",
	Description: "used to classify the intent of the user.",
}

type Classifier struct {
	*node.Prompt
	subchainMap map[string]node.Info
}

func NewClassifier(subchains ...node.Info) *Classifier {

	subchainMap := map[string]node.Info{}
	for _, s := range subchains {
		subchainMap[s.Name] = s
	}

	return &Classifier{
		subchainMap: subchainMap,
		Prompt: node.NewPrompt(
			ClassifierInfo,
			classifierTemplate,
			llm.OptionStop("\n"),
			// llm.OptionMaxTokens(100),
		),
	}
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
		in.
			WithKeyValue("subchains", n.subchainMap).
			WithKeyValue("original-user-input", in.Input()),
	)
}
