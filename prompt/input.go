package prompt

type Input interface {
	Input() string
	Get(key string) any
	Keys() map[string]any
}

type input struct {
	keys  map[string]any
	input string
}

func NewInput(in string) Input {
	return NewInputWithKeys(in, nil)
}

func NewInputWithKeys(in string, keys map[string]any) Input {
	return input{
		input: in,
		keys:  keys,
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
