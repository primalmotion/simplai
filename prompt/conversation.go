package prompt

import (
	"context"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

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

var ConversationDesc = node.Desc{
	Name:        "conversation",
	Description: "Used to have a general conversation with the user.",
}

type Conversation struct {
	*node.Prompt
	conversation *node.ChatMemory
}

func NewConversation() *Conversation {
	return &Conversation{
		Prompt: node.NewPrompt(
			ConversationDesc,
			conversationTemplate,
		),
	}
}

func (n *Conversation) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Conversation) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (n *Conversation) Execute(ctx context.Context, in node.Input) (string, error) {
	return n.Prompt.Execute(ctx, in.WithOptions(
		llm.OptionStop(
			in.Get("username").(string),
			in.Get("botname").(string),
			in.Get("system").(string),
		),
	))
}
