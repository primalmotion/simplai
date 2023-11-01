package node

import "git.sr.ht/~primalmotion/simplai/llm"

// Input represents the data passed to
// a Node.
type Input struct {
	keys       map[string]any
	input      string
	scratchpad string
	llmOptions []llm.Option
	debug      bool
}

// NewInput returns a new Input with the given string
// and llm.Options. Note that the llm.Option will be
// carried out until it reaches a LLM Node. After
// that they will be discarded.
func NewInput(in string, options ...llm.Option) Input {
	return Input{
		input:      in,
		keys:       map[string]any{},
		llmOptions: options,
	}
}

// WithInput returns a copy of the receiver
// replacing the previous input.
func (i Input) WithInput(in string) Input {
	i.input = in
	return i
}

// Input returns the current input string.
func (i Input) Input() string {
	return i.input
}

// Set returns a copy of the receiver with the given
// key and value added to it.
func (i Input) Set(k string, v any) Input {
	i.keys[k] = v
	return i
}

// Get returns the value of the given key.
func (i Input) Get(key string) any {
	return i.keys[key]
}

// WithLLMOptions returns a copy of the receiver with the given llm.Options
// added to it.
func (i Input) WithLLMOptions(options ...llm.Option) Input {
	i.llmOptions = append([]llm.Option{}, options...)
	return i
}

// LLMOptions returns the current llm.Options.
func (i Input) LLMOptions() []llm.Option {
	return i.llmOptions
}

// ResetLLMOptions returns a copy of the Input
// after removing all llm.Options. This is
// called by the LLM nodes.
func (i Input) ResetLLMOptions() Input {
	i.llmOptions = nil
	return i
}

// WithScratchpad returns a copy of the receiver with
// the given scratchpad added.
func (i Input) WithScratchpad(scratchpad string) Input {
	i.scratchpad = scratchpad
	return i
}

// Scratchpad returns the current scratchpad.
func (i Input) Scratchpad() string {
	return i.scratchpad
}

// WithDebug returns a copy of the receiver with the
// Debug flag activated. The Node will usually print
// detailed information about them when they deal with
// an input with Debug set to true.
func (i Input) WithDebug(debug bool) Input {
	i.debug = debug
	return i
}

// Debug returns the state of the debug mode.
func (i Input) Debug() bool {
	return i.debug
}
