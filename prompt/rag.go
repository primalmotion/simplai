package prompt

import (
	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/node"
)

// RagInfo is the node.Info for the Rag.
var RagInfo = node.Info{
	Name:        "retrieval augmented generation",
	Description: "use in context information to answer the query",
	Parameters:  "The query and context information retrieved",
}

const ragTemplate = `You are an assistant for question-answering tasks.
Use the following pieces of retrieved context only to answer the question.
If you don't find the answer, just say that you don't know.

Context:

{{ .Input }}

Question: {{ .Get "userquery" }}

Keep the answer concise. 
Only use the information provided above.

Answer:`

// A Rag is a prompt asking the LLM to
// perform rag operations.
type Rag struct {
	*node.Prompt
}

// NewRag returns a new *Rag.
func NewRag() *Rag {
	return &Rag{
		Prompt: node.NewPrompt(
			RagInfo,
			ragTemplate,
			engine.OptionStop("Question:", "Context:"),
		),
	}
}
