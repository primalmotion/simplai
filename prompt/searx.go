package prompt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/primalmotion/simplai/node"
)

// SearxSearchInfo is the node.Information for the SearxSearch prompt.
var SearxSearchInfo = node.Info{
	Name:        "search",
	Description: "use to access internet and find (search) current data about people, news, etc...",
	Parameters:  "the subject to search for",
}

const searxTemplate = `You must extract the information from the following
data. Write a short summary of about 2-3 sentences.

QUERY: {{ .Get "userquery" }}

RESULTS:

{{ .Input }}

SUMMARY:`

type searxTrimmedResponse struct {
	Results []struct {
		Content  string  `json:"content"`
		Title    string  `json:"title"`
		Category string  `json:"category"`
		URL      string  `json:"url"`
		Score    float64 `json:"score"`
	} `json:"results"`
}

// A SearxSearch is a prompt fetch data from intermet using Searx
// then summarizes the search results.
type SearxSearch struct {
	*node.Prompt
	client http.Client
	api    string
}

// NewSearxSearch returns a new *SearxSearch using
// the provided URL.
func NewSearxSearch(api string) *SearxSearch {
	client := http.Client{}
	return &SearxSearch{
		api:    api,
		client: client,
		Prompt: node.NewPrompt(
			SearxSearchInfo,
			searxTemplate,
		),
	}
}

// Execute implements the node.Node interface.
// It will make a search using the provided input, then massage the data
// and summarize them.
func (n *SearxSearch) Execute(ctx context.Context, in node.Input) (string, error) {

	query := in.Input()

	u, err := url.Parse(fmt.Sprintf("%s/search", n.api))
	if err != nil {
		return "", fmt.Errorf("unable to parse url: %w", err)
	}
	values := url.Values{
		"format": {"json"},
		"q":      {query},
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("unable to create http request: %w", err)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to perform http request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned an error: %s", resp.Status)
	}

	out := searxTrimmedResponse{}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&out); err != nil {
		return "", fmt.Errorf("unable to decode response: %w", err)
	}

	output := []string{}
	for _, entry := range out.Results {
		output = append(output, fmt.Sprintf(
			"- %s (score: %2f)\n%s",
			entry.Title,
			entry.Score,
			entry.Content,
		))
	}

	return n.Prompt.Execute(
		ctx,
		in.WithInput(strings.Join(output, "\n\n")).Set("userquery", query),
	)
}
