package node

import (
	"fmt"
	"strings"

	"git.sr.ht/~primalmotion/simplai/utils/trim"
)

type Conversation struct {
	*BaseNode
	botname  string
	username string
	history  []string
}

func NewConversation(botname string, username string) *Conversation {
	return &Conversation{
		BaseNode: New(),
		botname:  strings.ToUpper(botname),
		username: strings.ToUpper(username),
	}
}

func (c *Conversation) BotName() string {
	return c.botname
}

func (c *Conversation) UserName() string {
	return c.username
}

func (c *Conversation) History() []string {
	return append([]string{}, c.history...)
}

func (c *Conversation) AddUserMessage(content string) {
	c.history = append(c.history, fmt.Sprintf("%s: %s", c.username, content))
}

func (c *Conversation) AddBotMessage(content string) {
	c.history = append(c.history, fmt.Sprintf("%s: %s", c.botname, content))
}

func (n *Conversation) Name() string {
	return "prompt"
}

func (n *Conversation) WithPreHook(h PreHook) Node {
	n.BaseNode.WithPreHook(h)
	return n
}

func (n *Conversation) WithPostHook(h PostHook) Node {
	n.BaseNode.WithPostHook(h)
	return n
}

func (n *Conversation) Execute(input Input) (string, error) {

	input = input.
		WithKeyValue("botname", n.botname).
		WithKeyValue("username", n.username).
		WithKeyValue("history", n.History())

	n.AddUserMessage(input.Input())

	output, err := n.BaseNode.Execute(input)
	if err != nil {
		return "", err
	}

	n.AddBotMessage(trim.Output(output))

	return output, nil
}
