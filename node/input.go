package node

import "git.sr.ht/~primalmotion/simplai/llm"

type Input interface {
	Input() string
	Get(key string) any
	Keys() map[string]any
	Options() []llm.Option

	WithKeyValue(k string, v any) Input
}

type input struct {
	keys    map[string]any
	input   string
	options []llm.Option
}

func NewInput(in string, options ...llm.Option) Input {
	return &input{
		input:   in,
		keys:    map[string]any{},
		options: options,
	}
}

func (i *input) WithKeyValue(k string, v any) Input {
	i.keys[k] = v
	return i
}

func (i *input) Input() string {
	return i.input
}

func (i *input) Get(key string) any {
	return i.keys[key]
}

func (i *input) Keys() map[string]any {
	return i.keys
}

func (i *input) Options() []llm.Option {
	return i.options
}
