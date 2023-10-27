package classifier

import (
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/prompt/basic"
)

const tmpl = `classify the user input and map it to a known action name.
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

Remember: you can only user one of: {{ range $k, $v := .Keys}}{{$k}}, {{end}}

## PROCEED


INPUT: {{ .Input }}
ACTION:`

type classifier struct {
	basic.Formatter
}

func NewClassifier() prompt.Formatter {
	return &classifier{
		Formatter: basic.Formatter{
			Template: tmpl,
		},
	}
}
