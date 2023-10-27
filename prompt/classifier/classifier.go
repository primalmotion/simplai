package classifier

import (
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/prompt/basic"
)

const tmpl = `classify the given input to understand what the user wants. You
must write the name of one of the following actions the user wants you to do:

OUTPUT FORMAT:
{{ range $k, $v := .Keys }}
	Intent: {{ $v }}
	Action Name: {{ $k }}
{{ end }}

For example, if you have the action:

	Intent: invent a song.
	Action Name: [name of the action]

If this is what the user wants, you must write:

	Action Name: [name of the action]

- You must respect this format, no matter what.
- Do NOT write anything else. You MUST NOT invent an action name. Instead, you MUST write:

	Action Name: NONE

USER INPUT:

{{ .Input }}

Action Name:`

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
