package prompt

import (
	"context"
	"fmt"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

// ClassifierInfo is the node.Info for the Classifier.
var ClassifierInfo = node.Info{
	Name:        "classifier",
	Description: "used to classify the intent of the user.",
}

const classifierTemplate = `classify the user input and map it to a known tool name.
your output will be used by a machine and it needs to be perfect.
You must conform to the protocol described below at all times.

## PROTOCOL

You must identify which of the known tools is the best fulfill the user request.
They are described as:

	NAME: <name>
	USAGE: <intent>
	PARAMS: <information about parameters>

The tool NAME is what you must classify. You will use the USAGE to help you
choose a tool. You also need to pay close attention to the PARAMS, if any, in
order to complete the parameters.

You MUST write a valid JSON output, describing the tools to use in the form:

	{"name":"<tool-name>","input":"<required-tool-parameters"}

For example, given the tools:

	- NAME: code
	  USAGE: use to write some code in various programming language
	  PARAMS: the description of the desired code to produce

	- NAME: compose
	  USAGE: use to compose a song
	  PARAMS: the subject the song is about.

Here is some output example:

	INPUT: write a hello world program in python
	ACTION: {"name":"code","input":"hello world in python"}

	INPUT: compose a song about bananas
	ACTION: {"name":"compose","input":"bananas"}

	INPUT: jump over the bridge
	ACTION: {"name":"default","input":"jump over the bridge"}

If the input not explicitly map to any available tools, you MUST exactly write:

	ACTION: {"name":"default","input":"[original-input]"}

It is VERY IMPORTANT you remember that you MUST follow this protocol no matter
what, in all circumstances.

## AVAILABLE TOOLS
{{ range $v := .Get "tools" }}
	- NAME: {{ $v.Name }}
	  USAGE: {{ $v.Description }}
	  PARAMS: {{ $v.Parameters }}
{{ end }}

Remember: ACTION's name must only be one of:
{{- range $k, $v := .Get "tools" }}
- {{$v.Name}}
{{- end}}

Pay attention to the tools description if it details what the input should be.

{{- if .Scratchpad }}

## PREVIOUS OBSERVATION

The following is observations about one of your previous failed attempts.
make sure you take them into account when generating the response.

- {{ .Scratchpad }}
{{- end }}

## PROCEED

INPUT: {{ .Input }}
ACTION:`

// A Classifier is a prompt that will try to classify an input
// into using one of of the tools it knows. The tools are different
// node.Node that are identified by their node.Info.
type Classifier struct {
	*node.Prompt
	tools []node.Info
}

// NewClassifier returns a new *Classifier.
func NewClassifier(tools ...node.Info) *Classifier {
	return &Classifier{
		tools: tools,
		Prompt: node.NewPrompt(
			ClassifierInfo,
			classifierTemplate,
			llm.OptionStop("\n"),
		),
	}
}

// Execute implements the Node interface.
func (n *Classifier) Execute(ctx context.Context, in node.Input) (output string, err error) {

	if len(n.tools) == 0 {
		return fmt.Sprintf(`{"name":"default","input":"%s"}`, in.Input()), nil
	}

	return n.Prompt.Execute(ctx, in.Set("tools", n.tools))
}
