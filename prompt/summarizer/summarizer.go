package summarizer

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	readability "github.com/go-shiori/go-readability"
)

const tmpl = `You are extremely good at providing precise 3-4 sentence summary
of any text. You will now summarize the following input:

{{ .Input }}

SUMMARY:`

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

type Summarizer struct {
	*node.Prompt
}

func NewSummarizer() *Summarizer {
	return &Summarizer{
		Prompt: node.NewPrompt(tmpl),
	}
}

func (n *Summarizer) Name() string {
	return fmt.Sprintf("%s:summarizer", n.BaseNode.Name())
}

func (n *Summarizer) WithPreHook(h node.PreHook) node.Node {
	n.BaseNode.WithPreHook(h)
	return n
}

func (n *Summarizer) WithPostHook(h node.PostHook) node.Node {
	n.BaseNode.WithPostHook(h)
	return n
}

func (s *Summarizer) Execute(in prompt.Input) (string, error) {

	text := in.Input()

	if _, err := url.ParseRequestURI(text); err == nil {
		article, err := readability.FromURL(text, 30*time.Second)
		if err != nil {
			return "", fmt.Errorf("unable to load article: %w", err)
		}
		text = article.TextContent
	}

	text = standardizeSpaces(text)
	if len(text) > 2048 {
		text = text[:2048]
	}

	return s.Prompt.Execute(prompt.NewInput(text))
}
