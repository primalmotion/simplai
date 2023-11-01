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

type PromptError struct {
	Err   error
	Input Input
}

func NewPromptError(input Input, err error) PromptError {
	return PromptError{
		Err:   err,
		Input: input,
	}
}

func (e PromptError) Unwrap() error {
	return e.Err
}

func (e PromptError) Error() string {
	return fmt.Sprintf("prompt error: %s", e.Err)
}
