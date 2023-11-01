package node

import "fmt"

type Error struct {
	N   Node
	Err error
}

func NewError(n Node, template string, val ...any) Error {
	return Error{
		N:   n,
		Err: fmt.Errorf(template, val...),
	}
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Error() string {
	return fmt.Sprintf("node error: %s: %s", e.N.Info().Name, e.Err)
}

