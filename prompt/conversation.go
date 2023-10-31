package prompt

import (
	"context"
	"fmt"

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

type Conversation struct {
	*node.Prompt
	conversation *node.ChatMemory
}

func NewConversation(c *node.ChatMemory) *Conversation {
	return &Conversation{
		conversation: c,
		Prompt: node.NewPrompt(
			conversationTemplate,
			llm.OptionStop(
				fmt.Sprintf("%s", c.System()),
				fmt.Sprintf("%s", c.BotName()),
				fmt.Sprintf("%s", c.UserName()),
			),
		).
			WithName("conversation").
			WithDescription("Used to have a general conversation with the user.").(*node.Prompt),
	}
}

func (n *Conversation) WithName(name string) node.Node {
	n.Prompt.WithName(name)
	return n
}

func (n *Conversation) WithDescription(desc string) node.Node {
	n.Prompt.WithDescription(desc)
	return n
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
	return n.Prompt.Execute(ctx, in)
}