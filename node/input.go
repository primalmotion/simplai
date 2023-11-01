package node

import "git.sr.ht/~primalmotion/simplai/llm"

type Input struct {
	keys       map[string]any
	input      string
	llmOptions []llm.Option
	debug      bool
}

func NewInput(in string, options ...llm.Option) Input {
	return Input{
		input:      in,
		keys:       map[string]any{},
		llmOptions: options,
	}
}

func (i Input) WithKeyValue(k string, v any) Input {
	i.keys[k] = v
	return i
}

func (i Input) WithLLMOptions(options ...llm.Option) Input {
	i.llmOptions = append([]llm.Option{}, options...)
	return i
}

func (i Input) Input() string {
	return i.input
}

func (i Input) Get(key string) any {
	return i.keys[key]
}

func (i Input) LLMOptions() []llm.Option {
	return i.llmOptions
}

func (i Input) WithDebug(debug bool) Input {
	i.debug = debug
	return i
}

func (i Input) Debug() bool {
	return i.debug
}

func (i *Input) Derive(in string) Input {

	nkeys := make(map[string]any, len(i.keys))
	for k, v := range i.keys {
		nkeys[k] = v
	}

	return Input{
		input:      in,
		keys:       nkeys,
		debug:      i.debug,
		llmOptions: append([]llm.Option{}, i.llmOptions...),
	}
}
