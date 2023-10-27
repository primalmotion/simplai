package prompt

import "git.sr.ht/~primalmotion/simplai/llm"

type Input interface {
	Input() string
	Get(key string) any
	Keys() map[string]any
	Options() []llm.Option
}

type input struct {
	keys    map[string]any
	input   string
	options []llm.Option
}

func NewInput(in string, options ...llm.Option) Input {
	return NewInputWithKeys(in, nil, options...)
}

func NewInputWithKeys(in string, keys map[string]any, options ...llm.Option) Input {
	return input{
		input:   in,
		keys:    keys,
		options: options,
	}
}

func (i input) Input() string {
	return i.input
}

func (i input) Get(key string) any {
	if i.keys == nil {
		return nil
	}

	return i.keys[key]
}

func (i input) Keys() map[string]any {
	return i.keys
}

func (i input) Options() []llm.Option {
	return i.options
}
