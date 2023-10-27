package node

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/prompt"
)

type Prompt struct {
	*BaseNode
	template string
	options  []llm.Option
}

func NewPrompt(template string, options ...llm.Option) *Prompt {
	return &Prompt{
		template: template,
		options:  options,
		BaseNode: New(),
	}
}

func (n *Prompt) Name() string {
	return "prompt"
}

func (n *Prompt) Execute(input prompt.Input) (string, error) {

	tmpl, err := template.New("").Parse(n.template)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return n.BaseNode.Execute(
		prompt.NewInput(
			buf.String(),
			append(n.options, input.Options()...)...,
		),
	)
}
