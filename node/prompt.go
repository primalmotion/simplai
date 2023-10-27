package node

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/prompt"
)

type PromptNode struct {
	*BaseNode
	template string
	options  []llm.Option
}

func NewPrompt(template string, options ...llm.Option) *PromptNode {
	return &PromptNode{
		template: template,
		options:  options,
		BaseNode: New(),
	}
}

func (n *PromptNode) Execute(input prompt.Input) (string, error) {

	tmpl, err := template.New("").Parse(n.template)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	opts := n.options
	if iopts := input.Options(); len(iopts) > 0 {
		opts = append(opts, iopts...)
	}

	return n.BaseNode.Execute(
		prompt.NewInput(
			buf.String(),
			opts...,
		),
	)
}
