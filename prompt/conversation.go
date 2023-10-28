package prompt

import (
	"fmt"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
)

const conversationTemplate = `You name is {{ .Get "botname" }}. You are an IA
with high level of conversational skills. You are entairtaining and curious.
You are knowledgeable in programming, physics, artificial intelligence, biology
and philosophy.

You have a conversation with {{ .Get "username" }}. You will do your best to
continue the conversation. You will make use of the context provided below in
order to stay extremely coherent.

AI: Hello how may I help you?
{{- .Get "history" }}
{{ .Get "username" }}: {{ .Input }}
{{ .Get "botname" }}: `

type Conversation struct {
	*node.Prompt
	botname  string
	username string
	history  []string
}

func NewConversation(botname string, username string) *Conversation {
	botname = strings.ToUpper(botname)
	username = strings.ToUpper(username)
	return &Conversation{
		botname:  botname,
		username: username,
		Prompt: node.NewPrompt(
			conversationTemplate,
			llm.OptionStop(
				fmt.Sprintf("\n%s", botname),
				fmt.Sprintf("\n%s", username),
			),
		),
	}
}

func (n *Conversation) AddMessageToHistory(name string, content string) {
	name = strings.ToUpper(name)
	n.history = append(n.history, fmt.Sprintf("%s: %s", name, content))
}

func (n *Conversation) Name() string {
	return fmt.Sprintf("%s:conversation", n.Prompt.Name())
}

func (n *Conversation) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Conversation) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (n *Conversation) Execute(in node.Input) (string, error) {

	in = in.
		WithKeyValue("botname", n.botname).
		WithKeyValue("username", n.username).
		WithKeyValue("history", strings.Join(n.history, "\n"))

	n.AddMessageToHistory(n.username, in.Input())

	output, err := n.Prompt.Execute(in)
	if err != nil {
		return "", err
	}

	n.AddMessageToHistory(n.botname, output)

	return output, nil
}
