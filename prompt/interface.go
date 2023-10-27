package prompt

type Formatter interface {
	Format(Input) (string, error)
}

type Input interface {
	Input() string
	Get(key string) any
}

type input struct {
	keys  map[string]any
	input string
}

func NewInput(in string, keys map[string]any) Input {
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
