package node

import (
	"bytes"
	"context"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/llm"
	"github.com/Masterminds/sprig"
)

type Prompt struct {
	*BaseNode
	template string
	options  []llm.Option
}

func NewPrompt(info Info, template string, options ...llm.Option) *Prompt {
	return &Prompt{
		template: template,
		options:  options,
		BaseNode: New(info),
	}
}

func (n *Prompt) Options() []llm.Option {
	return append([]llm.Option{}, n.options...)
}

func (n *Prompt) Execute(ctx context.Context, input Input) (string, error) {

	tmpl, err := template.New("base").
		Funcs(sprig.FuncMap()).
		Parse(n.template)
	if err != nil {
		return "", NewError(n, "unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", NewError(n, "unable to execute template: %w", err)
	}

	if input.Debug() {
		LogNode(n, "2", buf.String())
	}

	return n.BaseNode.Execute(
		ctx,
		input.
			Derive(buf.String()).
			WithLLMOptions(append(n.options, input.LLMOptions()...)...),
	)
}
