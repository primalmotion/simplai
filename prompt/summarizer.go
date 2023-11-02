package prompt

import (
	"context"
	"fmt"
	"net/url"
	"time"

	readability "github.com/go-shiori/go-readability"
	"github.com/primalmotion/simplai/node"
	"github.com/primalmotion/simplai/utils/trim"
)

// SummarizerInfo is the node.Info for the Summarizer.
var SummarizerInfo = node.Info{
	Name:        "summarizer",
	Description: "use to summarize (resume, brief, shorten) some text, URL or document.",
	Parameters:  "either the full texAt to summarize or a valid URL. if the url schema is missing, assume https://",
}

const summarizerTemplate = `You are extremely good at providing precise 3-4 sentence summary
of any text. You will now summarize the following input:

{{ .Input }}

SUMMARY:`

// A Summarizer is a prompt that will try to summuarize the given
// input. If the input is a valid URL, the content of that URL will
// first be retrieved using readability, then will be summarized.
type Summarizer struct {
	*node.Prompt
}

// NewSummarizer returns a new *Summarizer.
func NewSummarizer() *Summarizer {
	return &Summarizer{
		Prompt: node.NewPrompt(
			SummarizerInfo,
			summarizerTemplate,
		),
	}
}

// Execute implements the node.Node interface.
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
	if len(text) > 4096 {
		text = text[:4096]
	}

	return s.Prompt.Execute(ctx, in.WithInput(text))
}
