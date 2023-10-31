package node

import "git.sr.ht/~primalmotion/simplai/llm"

type Input struct {
	keys    map[string]any
	input   string
	options []llm.Option
}

func NewInput(in string, options ...llm.Option) Input {
	return Input{
		input:   in,
		keys:    map[string]any{},
		options: options,
	}
}

func (i Input) WithKeyValue(k string, v any) Input {
	i.keys[k] = v
	return i
}

func (i Input) WithOptions(options ...llm.Option) Input {
	i.options = append(i.options, options...)
	return i
}

func (i Input) Input() string {
	return i.input
}

func (i Input) Get(key string) any {
	return i.keys[key]
}

func (i Input) Keys() map[string]any {
	return i.keys
}

func (i Input) Options() []llm.Option {
	return i.options
}

func (i *Input) Derive(in string) Input {

	nkeys := make(map[string]any, len(i.keys))
	for k, v := range i.keys {
		nkeys[k] = v
	}

	return Input{
		input:   in,
		keys:    nkeys,
		options: append([]llm.Option{}, i.options...),
	}
}
