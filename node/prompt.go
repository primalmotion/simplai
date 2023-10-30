package node

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/llm"
	"github.com/Masterminds/sprig"
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

func (n *Prompt) WithPreHook(h PreHook) Node {
	n.BaseNode.WithPreHook(h)
	return n
}

func (n *Prompt) WithPostHook(h PostHook) Node {
	n.BaseNode.WithPostHook(h)
	return n
}

func (n *Prompt) Execute(input Input) (string, error) {

	tmpl, err := template.New("").
		Funcs(sprig.FuncMap()).
		Parse(n.template)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return n.BaseNode.Execute(
		NewInput(
			buf.String(),
			append(n.options, input.Options()...)...,
		),
	)
}
