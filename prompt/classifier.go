package prompt

import (
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

Example:

	code: write code using various programming languages.
	compose: write or compose a song.

Given a user input, you must find which action is describing the best the
user's intent. For example, if the user input is:

	INPUT: write a hello world program in python
	ACTION: code
	INPUT: compose a song about bananas
	ACTION: compose

If a sequence of multiple actions is needed, you must write them in order:

	INPUT: summarize the latest news and post them on the my blog
	ACTION: search, summarize, post

If the input not explicitely map to any known actions, you must write:

	ACTION: none

It is VERY IMPORTANT you understand that you MUST follow this protocol at all
costs, no matter what, and at all circumstances.

## KNOWN ACTIONS
{{ range $k, $v := .Keys }}
	{{ $k }}: {{ $v -}}
{{ end }}

Remember: you can write one of: {{ range $k, $v := .Keys}}{{$k}}, {{end}}

## PROCEED

INPUT: {{ .Input }}
ACTION:`

type Classifier struct {
	*node.Prompt
}

func NewClassifier() *Classifier {
	return &Classifier{
		Prompt: node.NewPrompt(
			classifierTemplate,
			llm.OptionStop("\n", " ", ","),
			llm.OptionMaxTokens(10),
		),
	}
}

func (n *Classifier) Name() string {
	return fmt.Sprintf("%s:classifier", n.Prompt.Name())
}

func (n *Classifier) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Classifier) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}
