package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~primalmotion/simplai/node"
)

// SearxInfo is the node.Information for the SearxSearch prompt.
var SearxInfo = node.Info{
	Name:        "search",
	Description: "use to access internet and find (search) current data about people, news, etc...",
	Parameters:  "the subject to search for",
}

type searxTrimmedResponse struct {
	Results []struct {
		Content  string  `json:"content"`
		Title    string  `json:"title"`
		Category string  `json:"category"`
		URL      string  `json:"url"`
		Score    float64 `json:"score"`
	} `json:"results"`
}

// A Searx is a tool to fetch data from intermet using Searx
// then massage the data.
type Searx struct {
	*node.BaseNode
	client http.Client
	api    string
}

// NewSearx returns a new *Searx using
// the provided URL.
func NewSearx(api string) *Searx {
	client := http.Client{}
	return &Searx{
		api:    api,
		client: client,
		BaseNode: node.New(
			SearxInfo,
		),
	}
}

// Execute implements the node.Node interface.
// It will make a search using the provided input, then
// massage the output and set the input.Set("userquery")
func (n *Searx) Execute(ctx context.Context, in node.Input) (string, error) {

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

	return n.BaseNode.Execute(
		ctx,
		in.WithInput(strings.Join(output, "\n\n")).Set("userquery", query),
	)
}
