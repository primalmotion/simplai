package prompt

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/utils/trim"
	readability "github.com/go-shiori/go-readability"
)

const summarizerTemplate = `You are extremely good at providing precise 3-4 sentence summary
of any text. You will now summarize the following input:

{{ .Input }}

SUMMARY:`

type Summarizer struct {
	*node.Prompt
}

func NewSummarizer() *Summarizer {
	return &Summarizer{
		Prompt: node.NewPrompt(summarizerTemplate).
			WithName("summarizer").
			WithDescription("summarize some text, URL or document.").(*node.Prompt),
	}
}

func (n *Summarizer) WithName(name string) node.Node {
	n.Prompt.WithName(name)
	return n
}

func (n *Summarizer) WithDescription(desc string) node.Node {
	n.Prompt.WithDescription(desc)
	return n
}

func (n *Summarizer) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *Summarizer) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (s *Summarizer) Execute(ctx context.Context, in node.Input) (string, error) {

	text := in.Input()

	if _, err := url.ParseRequestURI(text); err == nil {
		article, err := readability.FromURL(text, 30*time.Second)
		if err != nil {
			return "", fmt.Errorf("unable to load article: %w", err)
		}
		text = article.TextContent
	}

	text = trim.Output(text)
	if len(text) > 2048 {
		text = text[:2048]
	}

	return s.Prompt.Execute(ctx, in.Derive(text))
}
