package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/primalmotion/simplai/utils/trim"
)

// ChatMemory holds information
// about chat history. This will
// be rewritten as this does not
// belong here.
type ChatMemory struct {
	*BaseNode
	history   *[]string
	system    string
	botname   string
	username  string
	separator string
}

// NewChatMemory returns a new ChatMemory with the provided information.
func NewChatMemory(info Info, system string, botname string, username string) *ChatMemory {
	return &ChatMemory{
		BaseNode:  New(info),
		system:    strings.ToLower(system),
		botname:   strings.ToLower(botname),
		username:  strings.ToLower(username),
		separator: "\n",
	}
}

// WithStorage sets the storage backend. Yeah it's an array..
// I told you this will be rewritten.
func (c *ChatMemory) WithStorage(storage *[]string) *ChatMemory {
	c.history = storage
	return c
}

// BotName returns the current botname.
func (c *ChatMemory) BotName() string {
	return c.botname
}

// UserName returns the current username.
func (c *ChatMemory) UserName() string {
	return c.username
}

// System returns the current system name.
func (c *ChatMemory) System() string {
	return c.system
}

// History returns the current history.
func (c *ChatMemory) History() []string {
	return append([]string{}, *c.history...)
}

// AddUserMessage add a new user message to the history
func (c *ChatMemory) AddUserMessage(content string) {
	*c.history = append(
		*c.history,
		fmt.Sprintf("%s%s%s", c.username, c.separator, content),
	)
}

// AddBotMessage add a new bot message to the history
func (c *ChatMemory) AddBotMessage(content string) {
	*c.history = append(
		*c.history,
		fmt.Sprintf("%s%s%s", c.botname, c.separator, content),
	)
}

// Execute implements the Node interface.
func (c *ChatMemory) Execute(ctx context.Context, input Input) (string, error) {

	if input.Input() == "flush" {
		*c.history = []string{}
		return "memory flushed", nil
	}

	input = input.
		Set("system", c.system).
		Set("botname", c.botname).
		Set("username", c.username).
		Set("history", c.History())

	output, err := c.BaseNode.Execute(ctx, input)
	if err != nil {
		return "", err
	}

	c.AddUserMessage(input.Input())
	c.AddBotMessage(trim.Output(output))

	return output, nil
}
