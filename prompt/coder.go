package prompt

import (
	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

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

type Coder struct {
	*node.Prompt
}

func NewCoder() *Coder {
	return &Coder{
		Prompt: node.NewPrompt(
			CoderInfo,
			coderTemplate,
			llm.OptionStop("<|EOF|>", "\nINPUT"),
		),
	}
}