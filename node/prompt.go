package node

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/prompt"
)

type PromptNode struct {
	*BaseNode
	Template string
	Stop     []string
}

func NewPrompt(template string, stopWords []string) *PromptNode {
	return &PromptNode{
		Template: template,
		Stop:     stopWords,
		BaseNode: New(),
	}
}

func (s *PromptNode) Execute(input prompt.Input) (string, error) {

	tmpl, err := template.New("").Parse(s.Template)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return s.BaseNode.Execute(prompt.NewInput(buf.String(), s.Stop...))
}
