package prompt

import (
	"context"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

// ConversationInfo is the node.Info for Conversation.
var ConversationInfo = node.Info{
	Name:        "conversation",
	Description: "used to have a general conversation with the user",
	Parameters:  "the user INPUT, as is. You must not modify or summarize it",
}

const conversationTemplate = `{{ .Get "system" }}
You are an AI with high conversational skills. You are entertaining and
curious. You are knowledgeable in programming, physics, artificial
intelligence, biology and philosophy. You will do your best to continue the
conversation. You will make use of the context provided below in order to stay
extremely coherent. If you are not sure, just tell that you don't know.

{{ .Get "botname" }}
Hello how may I help you?
{{ .Get "history" | join "\n" }}
{{ .Get "username" }}
{{ .Input }}
{{ .Get "botname" }}
`

// A Conversation is a prompt that can be used
// to have a generic conversation with the LLM.
type Conversation struct {
	*node.Prompt
}

// NewConversation returns a new *Conversation.
func NewConversation() *Conversation {
	return &Conversation{
		Prompt: node.NewPrompt(
			ConversationInfo,
			conversationTemplate,
			llm.OptionTemperature(0.8),
		),
	}
}

// Execute implements the node.Node interface. It will
// inject into the input the llm.OptionStop needed.
func (n *Conversation) Execute(ctx context.Context, in node.Input) (string, error) {
	return n.Prompt.Execute(ctx, in.WithLLMOptions(
		llm.OptionStop(
			in.Get("username").(string),
			in.Get("botname").(string),
			in.Get("system").(string),
		),
	))
}
