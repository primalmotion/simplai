package prompt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"git.sr.ht/~primalmotion/simplai/node"
)

const searxTemplate = `You must read, understand and summarize
The following JSON data coming from the API of a search engine
called searx. You must extract the data and summarize the results.

The search the made was: {{ .Get "userquery" }}

API DATA:

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

type SearxSearch struct {
	*node.Prompt
	conversation *node.ChatMemory
	client       http.Client
	api          string
}

func NewSearxSearch(conversation *node.ChatMemory, api string) *SearxSearch {
	client := http.Client{}
	return &SearxSearch{
		api:          api,
		client:       client,
		conversation: conversation,
		Prompt:       node.NewPrompt(searxTemplate),
	}
}

func (n *SearxSearch) Name() string {
	return fmt.Sprintf("%s:searxsearch", n.Prompt.Name())
}

func (n *SearxSearch) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *SearxSearch) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}

func (n *SearxSearch) Execute(in node.Input) (string, error) {

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

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("unable to reencode json data: %w", err)
	}

	output, err := n.Prompt.Execute(
		node.NewInput(string(data)).
			WithKeyValue("userquery", query),
	)
	if err != nil {
		return "", fmt.Errorf("unable to execute query: %w", err)
	}

	return output, nil
}
