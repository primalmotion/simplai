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
	return fmt.Sprintf("[%s]: %s", e.N.Info().Name, e.Err)
}

func (e Error) Is(err error) bool {
	_, ok := err.(Error)
	return ok
}

type PromptError struct {
	Err        error
	Scratchpad string
}

func NewPromptError(scratchpad string, err error) PromptError {
	return PromptError{
		Err:        err,
		Scratchpad: scratchpad,
	}
}

func (e PromptError) Unwrap() error {
	return e.Err
}

func (e PromptError) Error() string {
	return fmt.Sprintf("prompt-error: %s", e.Err)
}

func (e PromptError) Is(err error) bool {
	_, ok := err.(PromptError)
	return ok
}
