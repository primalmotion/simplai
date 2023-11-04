package node

import (
	"context"

	"github.com/primalmotion/simplai/engine"
)

// LLM is a node responsible for running its
// input into an inference engine.
type LLM struct {
	engine engine.LLM
	*BaseNode
	options []engine.LLMOption
}

// NewLLM returns a new LLM using the given engine.
func NewLLM(info Info, engine engine.LLM, options ...engine.LLMOption) *LLM {
	return &LLM{
		BaseNode: New(info),
		engine:   engine,
		options:  options,
	}
}

// Execute implements the Node interface.
func (n *LLM) Execute(ctx context.Context, input Input) (string, error) {

	opts := append(n.options, input.LLMOptions()...)
	if input.Debug() {
		opts = append(opts, engine.OptionDebug(true))
	}

	output, err := n.engine.Infer(ctx, input.Input(), opts...)
	if err != nil {
		return "", NewError(n, "unable to run llm inference: %w", err)
	}

	if input.Debug() {
		LogNode(n, "4", output)
	}

	return n.BaseNode.Execute(
		ctx,
		input.
			WithInput(output).
			ResetLLMOptions(),
	)
}
