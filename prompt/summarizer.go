package prompt

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.sr.ht/~primalmotion/simplai/node"
	readability "github.com/go-shiori/go-readability"
)

const summarizerTemplate = `You are extremely good at providing precise 3-4 sentence summary
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
		Prompt: node.NewPrompt(summarizerTemplate),
	}
}

func (n *Summarizer) Name() string {
	return fmt.Sprintf("%s:summarizer", n.Prompt.Name())
}

func (n *Summarizer) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Summarizer) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (s *Summarizer) Execute(in node.Input) (string, error) {

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

	return s.Prompt.Execute(node.NewInput(text))
}
