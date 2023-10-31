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

var SummarizerDesc = node.Info{
	Name:        "summarizer",
	Description: "summarize some text, URL or document.",
	Parameters:  "either the full text to summarize or a valid URL",
}

const summarizerTemplate = `You are extremely good at providing precise 3-4 sentence summary
of any text. You will now summarize the following input:

{{ .Input }}

SUMMARY:`

type Summarizer struct {
	*node.Prompt
}

func NewSummarizer() *Summarizer {
	return &Summarizer{
		Prompt: node.NewPrompt(
			SummarizerDesc,
			summarizerTemplate,
		),
	}
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
