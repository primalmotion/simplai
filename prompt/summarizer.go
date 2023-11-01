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

var SummarizerInfo = node.Info{
	Name:        "summarizer",
	Description: "use to summarize (resume, brief, shorten) some text, URL or document.",
	Parameters:  "either the full texAt to summarize or a valid URL. if the url schema is missing, assume https://",
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
			SummarizerInfo,
			summarizerTemplate,
		),
	}
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
