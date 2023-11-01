package node

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/llm"
	"github.com/Masterminds/sprig"
)

type Prompt struct {
	*BaseNode
	template   string
	options    []llm.Option
	maxRetries int
}

func NewPrompt(info Info, template string, options ...llm.Option) *Prompt {
	return &Prompt{
		template:   template,
		options:    options,
		maxRetries: 3,
		BaseNode:   New(info),
	}
}

func (n *Prompt) Options() []llm.Option {
	return append([]llm.Option{}, n.options...)
}

func (n *Prompt) WithMaxRetries(maxRetries int) *Prompt {
	n.maxRetries = maxRetries
	return n
}

func (n *Prompt) Execute(ctx context.Context, input Input) (output string, err error) {

	tmpl, err := template.New("base").
		Funcs(sprig.FuncMap()).
		Parse(n.template)
	if err != nil {
		return "", NewError(n, "unable to parse template: %w", err)
	}

	for i := 0; i < n.maxRetries; i++ {

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, input); err != nil {
			return "", NewError(n, "unable to execute template: %w", err)
		}

		if input.Debug() {
			LogNode(n, "2", buf.String())
		}

		output, err = n.BaseNode.Execute(
			ctx,
			input.Derive(buf.String()).WithLLMOptions(
				append(
					n.options,
					input.LLMOptions()...,
				)...,
			),
		)

		if err == nil {
			return output, nil
		}

		var promptErr PromptError
		if !errors.As(err, &promptErr) {
			return "", err
		}

		if input.Debug() {
			LogNode(n, "3", fmt.Sprintf(
				"Got a PromptError\n\nScratchpad: %s\nRemaining: %d\nError: %s",
				promptErr.Scratchpad,
				n.maxRetries-i,
				err,
			))
		}

		input = input.WithScratchpad(promptErr.Scratchpad)
	}

	return "", err
}
