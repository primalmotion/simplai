package node

import (
	"context"
	"fmt"
	"strings"

	"git.sr.ht/~primalmotion/simplai/utils/trim"
)

type ChatMemory struct {
	*BaseNode
	history   *[]string
	system    string
	botname   string
	username  string
	separator string
}

func NewChatMemory(desc Info, system string, botname string, username string) *ChatMemory {
	return &ChatMemory{
		BaseNode:  New(desc),
		system:    strings.ToLower(system),
		botname:   strings.ToLower(botname),
		username:  strings.ToLower(username),
		separator: "\n",
	}
}

func (c *ChatMemory) WithStorage(storage *[]string) *ChatMemory {
	c.history = storage
	return c
}

func (c *ChatMemory) BotName() string {
	return c.botname
}

func (c *ChatMemory) UserName() string {
	return c.username
}

func (c *ChatMemory) System() string {
	return c.system
}

func (c *ChatMemory) History() []string {
	return append([]string{}, *c.history...)
}

func (c *ChatMemory) AddUserMessage(content string) {
	*c.history = append(
		*c.history,
		fmt.Sprintf("%s%s%s", c.username, c.separator, content),
	)
}

func (c *ChatMemory) AddBotMessage(content string) {
	*c.history = append(
		*c.history,
		fmt.Sprintf("%s%s%s", c.botname, c.separator, content),
	)
}

func (c *ChatMemory) Execute(ctx context.Context, input Input) (string, error) {

	if input.Input() == "flush" {
		*c.history = []string{}
		return "memory flushed", nil
	}

	input = input.
		WithKeyValue("system", c.system).
		WithKeyValue("botname", c.botname).
		WithKeyValue("username", c.username).
		WithKeyValue("history", c.History())

	output, err := c.BaseNode.Execute(ctx, input)
	if err != nil {
		return "", err
	}

	c.AddUserMessage(input.Input())
	c.AddBotMessage(trim.Output(output))

	return output, nil
}
