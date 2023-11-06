package tool

import (
	"context"
	"fmt"
	"strings"

	"github.com/primalmotion/simplai/node"
	"github.com/primalmotion/simplai/vectorstore"
)

// RetrieverInfo is the node.Information for retriever prompt.
var RetrieverInfo = node.Info{
	Name:        "retriever",
	Description: "use to retrieve contextual information related to the query",
	Parameters:  "the query to process",
}

// Retriever is a tool to retrieve data from vectorstore related to the query.
type Retriever struct {
	*node.BaseNode
	store vectorstore.VectorStore
	topk  int
}

// NewRetriever returns a new retriever.
func NewRetriever(store vectorstore.VectorStore, topk int) *Retriever {
	return &Retriever{
		BaseNode: node.New(RetrieverInfo),
		store:    store,
		topk:     topk,
	}
}

// Execute implements the node.Node interface.
// It will make a search using the provided input, then
// massage the output and set the the topk as input.Set("results") as an array
// if we need to post process them with a reranker for instance.
func (n *Retriever) Execute(ctx context.Context, in node.Input) (string, error) {

	query := in.Input()

	docs, err := n.store.SimilaritySearch(ctx, query, n.topk)
	if err != nil {
		return "", err
	}

	output := []string{}
	for _, entry := range docs {
		output = append(output, fmt.Sprintf(
			"- %s (score: %.2f)\n%s",
			entry.ID,
			entry.Distance,
			entry.Content,
		))
	}

	return n.BaseNode.Execute(
		ctx,
		in.WithInput(strings.Join(output, "\n\n")).Set("userquery", query).Set("results", docs),
	)
}
