package prompt

import (
	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/node"
)

// CoderInfo is the node.Info for the Coder.
var CoderInfo = node.Info{
	Name:        "coder",
	Description: "use to write some code, in various programming language",
	Parameters:  "The detailed summary of the code to write",
}

const coderTemplate = `You are a skilled programmer able to write very good and
efficient programs and code snippets. You can code in any language, but you are
particularly proficient in Go, Python, Bash and Javascript.

If you are are to code in Java, just respond "I don't code in that language and
you shouldn't too".

When you are finished, you MUST write a new single line containing only "<|EOF|>".

# EXAMPLE

INPUT: write a hello world in bash
CODE: #!/bin/bash
echo "hello world"
<|EOF|>

# PROCEED

INPUT: {{ .Input }}
CODE:`

// A Coder is a prompt asking the LLM to
// perform various coding operations.
type Coder struct {
	*node.Prompt
}

// NewCoder returns a new *Coder.
func NewCoder() *Coder {
	return &Coder{
		Prompt: node.NewPrompt(
			CoderInfo,
			coderTemplate,
			engine.OptionStop("<|EOF|>", "\nINPUT"),
		),
	}
}
