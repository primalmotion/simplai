package tool

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-shiori/go-readability"
	"github.com/primalmotion/simplai/node"
)

// ScraperInfo is the node.Information for the scraper tool.
var ScraperInfo = node.Info{
	Name:        "scraper",
	Description: "use to scrap and convert an url into a list of documents",
	Parameters:  "the url to scrap",
}

// A scraper is a tool to scrape an url and then massage the data
// into a list of vectorstore.Document
type Scraper struct {
	*node.BaseNode
	client http.Client
}

// NewSearx returns a new *Searx using
// the provided URL.
func NewScraper() *Scraper {
	client := http.Client{}
	return &Scraper{
		client: client,
		BaseNode: node.New(
			ScraperInfo,
		),
	}
}

// Execute implements the node.Node interface.
// It will make a scrap using the provided input, then
// massage the output and set the input.Set("userquery")
// output is also converted as []vectorstore.Document and set
// as the input.Set("documents") for tool chaining.
func (n *Scraper) Execute(ctx context.Context, in node.Input) (string, error) {

	rawurl := in.Input()

	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("unable to parse url %s: %w", rawurl, err)
	}

	buffer := bytes.NewBuffer(nil)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), buffer)

	if err != nil {
		return "", fmt.Errorf("unable to prepare request: %w", err)
	}

	resp, err := n.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("unable to send request: %w", err)
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		content, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server was unable to process the request: %s\n\n%s", resp.Status, content)
	}

	article, err := readability.FromReader(resp.Body, parsedUrl)
	if err != nil {
		return "", fmt.Errorf("unable to parse %s: %w", parsedUrl.String(), err)
	}

	// TODO
	fmt.Println(article)
	var output []string
	var query string
	var docs string

	// Need to split

	// output := []string{}
	// docs := []vectorstore.Document{}
	// for _, entry := range out.Results {
	// 	docs = append(docs, vectorstore.Document{
	// 		Metadata: map[string]any{
	// 			"URL":      entry.URL,
	// 			"Category": entry.Category,
	// 			"Score":    entry.Score,
	// 		},
	// 		ID:      entry.Title,
	// 		Content: entry.Content,
	// 	})
	// 	output = append(output, fmt.Sprintf(
	// 		"- %s (score: %2f)\n%s",
	// 		entry.Title,
	// 		entry.Score,
	// 		entry.Content,
	// 	))
	// }

	return n.BaseNode.Execute(
		ctx,
		in.WithInput(strings.Join(output, "\n\n")).Set("userquery", query).Set("documents", docs),
	)
}
