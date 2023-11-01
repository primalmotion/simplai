package node

import "fmt"

// An Error represents an error returned
// by a node.
type Error struct {
	N   Node
	Err error
}

// NewError retrurns a new Error.
func NewError(n Node, template string, val ...any) Error {
	return Error{
		N:   n,
		Err: fmt.Errorf(template, val...),
	}
}

// Unwrap implements the error interface.
func (e Error) Unwrap() error {
	return e.Err
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("[%s]: %s", e.N.Info().Name, e.Err)
}

// Is implements the error interface.
func (e Error) Is(err error) bool {
	_, ok := err.(Error)
	return ok
}

// PromptError represents an error in the generated
// text generation. The scratchpad will be used by
// the Prompt to fill the input.Scratchpad before retrying.
type PromptError struct {
	Err        error
	Scratchpad string
}

// NewPromptError returns a PromptError.
func NewPromptError(scratchpad string, err error) PromptError {
	return PromptError{
		Err:        err,
		Scratchpad: scratchpad,
	}
}

// Unwrap implements the error interface.
func (e PromptError) Unwrap() error {
	return e.Err
}

// Error implements the error interface.
func (e PromptError) Error() string {
	return fmt.Sprintf("prompt-error: %s", e.Err)
}

// Is implements the error interface.
func (e PromptError) Is(err error) bool {
	_, ok := err.(PromptError)
	return ok
}
