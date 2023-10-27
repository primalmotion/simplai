package prompt

type Input interface {
	Input() string
	Get(key string) any
	Keys() map[string]any
	StopWords() []string
}

type input struct {
	keys  map[string]any
	input string
	stop  []string
}

func NewInput(in string, stop ...string) Input {
	return NewInputWithKeys(in, nil, stop...)
}

func NewInputWithKeys(in string, keys map[string]any, stop ...string) Input {
	return input{
		input: in,
		keys:  keys,
		stop:  stop,
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

func (i input) StopWords() []string {
	return i.stop
}
