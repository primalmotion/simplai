package basic

import (
	"bytes"
	"fmt"
	"text/template"

	"git.sr.ht/~primalmotion/simplai/prompt"
)

type Formatter struct {
	Template string
	Stop     []string
}

func (s *Formatter) Format(input prompt.Input) (string, error) {

	tmpl, err := template.New("").Parse(s.Template)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return buf.String(), nil
}

func (s *Formatter) StopWords() []string {
	return s.Stop
}
