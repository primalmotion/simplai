package websummarizer

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
)

const template = `You are extremely good at providing precise summary of any text.
You will now summarize the following input:

%s

SUMMARY:`

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

type WebSummarizer struct {
}

func (s *WebSummarizer) Format(address string) (string, error) {

	if _, err := url.Parse(address); err != nil {
		return "", fmt.Errorf("The given url is not valid: %w", err)
	}

	article, err := readability.FromURL(address, 30*time.Second)
	if err != nil {
		return "", fmt.Errorf("unable to load article: %w", err)
	}

	return fmt.Sprintf(template, standardizeSpaces(article.TextContent)[:2048]), nil
}
