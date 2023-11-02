package node

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/primalmotion/simplai/llm"
)

// A Prompt is a Node that is responsible
// for generating a prompt from a template using
// the informations contained in the given input.
type Prompt struct {
	*BaseNode
	template   string
	options    []llm.Option
	maxRetries int
}

// NewPrompt returns a new *Prompt with the given template and ll.Options.
// The llm.Options will always be appended to the input, before the input's
// own LLMOptions. So input's options will take precedence.
func NewPrompt(info Info, template string, options ...llm.Option) *Prompt {
	return &Prompt{
		template:   template,
		options:    options,
		maxRetries: 3,
		BaseNode:   New(info),
	}
}

// Options returns the current llm.Options.
func (n *Prompt) Options() []llm.Option {
	return append([]llm.Option{}, n.options...)
}

// WithMaxRetries sets the maximum number of retries the prompt
// is willing to make when the execution returns a PromptError.
func (n *Prompt) WithMaxRetries(maxRetries int) *Prompt {
	n.maxRetries = maxRetries
	return n
}

// Execute implements the Node interface.
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
			input.WithInput(buf.String()).WithLLMOptions(
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
