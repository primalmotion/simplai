package node

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/prompt"
)

type PromptNode struct {
	*BaseNode
	template string
	stop     []string
}

func NewPrompt(template string, stopWords []string) *PromptNode {
	return &PromptNode{
		template: template,
		stop:     stopWords,
		BaseNode: New(),
	}
}

func (s *PromptNode) Execute(input prompt.Input) (string, error) {

	tmpl, err := template.New("").Parse(s.template)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return s.BaseNode.Execute(
		prompt.NewInput(
			buf.String(),
			llm.OptionStop(s.stop...),
		),
	)
}
