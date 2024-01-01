package tool

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/node"
	"github.com/primalmotion/simplai/utils/reorder"
	"github.com/primalmotion/simplai/vectorstore"
)

// RerankerInfo is the node.Information for reranker prompt.
var RerankerInfo = node.Info{
	Name:        "reranker",
	Description: "use to re-rank a set of information for a given query",
	Parameters:  "the query and set of information to process",
}

// Reranker is a tool to rerank passes against a query.
type Reranker struct {
	*node.BaseNode
	reranker engine.Reranker
	topk     int
}

// NewReranker returns a new reranker.
func NewReranker(reranker engine.Reranker, topk int) *Reranker {
	return &Reranker{
		BaseNode: node.New(RerankerInfo),
		reranker: reranker,
		topk:     topk,
	}
}

// Execute implements the node.Node interface.
// It will rerank the messages with the initial query
func (n *Reranker) Execute(ctx context.Context, in node.Input) (string, error) {

	query := in.Get("userquery").(string)
	docs := in.Get("documents").([]vectorstore.Document)

	passages := make([]string, len(docs))

	for _, d := range docs {
		passages = append(passages, d.Content)
	}

	// Compute our rerank
	scores := make(map[float64][]vectorstore.Document)

	sl, err := n.reranker.Rerank(ctx, query, passages)
	if err != nil {
		return "", fmt.Errorf("unable to rerank documents: %w", err)
	}

	if len(sl) != len(passages) {
		return "", fmt.Errorf("unable to renrank documents the scores doesnt match the number of documents")
	}

	// note that is totaly possible to have same rerank scores
	for i, d := range docs {
		scores[sl[i]] = append(scores[sl[i]], d)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(sl)))

	output := []string{}

	for score, docs := range scores {
		for _, d := range docs {
			if len(output) == n.topk {
				break
			}
			output = append(output, fmt.Sprintf("- reranking score: %.2f\n  similarity score: %.2f\n  content:%s", score, d.Distance, d.Content))
		}
	}

	return n.BaseNode.Execute(
		ctx,
		in.WithInput(strings.Join(reorder.Distribute(output), "\n\n")).Set("userquery", query),
	)
}
