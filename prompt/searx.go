package prompt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~primalmotion/simplai/node"
)

const searxTemplate = `You must extract the information from the following
data. Write a short summary of about 2-3 sentences.

QUERY: {{ .Get "userquery" }}

RESULTS:

{{ .Input }}

SUMMARY:`

var SearxSearchDesc = node.Desc{
	Name:        "search",
	Description: "used to summarize some text, URL or document.",
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

type SearxSearch struct {
	*node.Prompt
	client http.Client
	api    string
}

func NewSearxSearch(api string) *SearxSearch {
	client := http.Client{}
	return &SearxSearch{
		api:    api,
		client: client,
		Prompt: node.NewPrompt(
			SearxSearchDesc,
			searxTemplate,
		),
	}
}

func (n *SearxSearch) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *SearxSearch) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

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
		in.Derive(strings.Join(output, "\n\n")).WithKeyValue("userquery", query),
	)
}
