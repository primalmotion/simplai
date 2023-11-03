package chromadb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	client *http.Client
	url    string
}

// New Return a new openAIAPI client.
func New(api string) *Client {

	client := &http.Client{}
	return &Client{
		url:    api,
		client: client,
	}
}

// CreateCollection creates a new collection. It will return it's ID.
func (v *Client) CreateCollection(ctx context.Context, req CollectionCreate) (CollectionResult, error) {

	res := CollectionResult{}
	err := v.send(
		ctx,
		"/api/v1/collections",
		http.MethodPost,
		&req,
		&res,
	)

	return res, err
}

// Add adds new embeddings.
func (v *Client) Add(ctx context.Context, collectionID string, req EmbeddingAdd) error {

	return v.send(
		ctx,
		fmt.Sprintf("/api/v1/collections/%s/add", collectionID),
		http.MethodPost,
		&req,
		nil,
	)
}

// Update update existing embeddings.
func (v *Client) Update(ctx context.Context, collectionID string, req EmbeddingUpdate) error {

	return v.send(
		ctx,
		fmt.Sprintf("/api/v1/collections/%s/update", collectionID),
		http.MethodPost,
		&req,
		nil,
	)
}

// Upsert updates existing embeddings if they exists or creates them if they
// don't.
func (v *Client) Upsert(ctx context.Context, collectionID string, req EmbeddingUpdate) error {

	return v.send(
		ctx,
		fmt.Sprintf("/api/v1/collections/%s/upsert", collectionID),
		http.MethodPost,
		&req,
		nil,
	)
}

// Delete deletes existing embeddings.
func (v *Client) Delete(ctx context.Context, collectionID string, req EmbeddingDelete) error {

	return v.send(
		ctx,
		fmt.Sprintf("/api/v1/collections/%s/delete", collectionID),
		http.MethodPost,
		&req,
		nil,
	)
}

func (v *Client) Get(ctx context.Context, collectionID string, req EmbeddingGet) (GetResult, error) {

	if len(req.Include) == 0 {
		req.Include = []Include{
			IncludeMetadatas,
			IncludeDocuments,
		}
	}

	res := GetResult{}
	err := v.send(
		ctx,
		fmt.Sprintf("/api/v1/collections/%s/get", collectionID),
		http.MethodPost,
		&req,
		&res,
	)
	return res, err
}

// Query queries for embeddings.
func (v *Client) Query(ctx context.Context, collectionID string, req EmbeddingQuery) (QueryResult, error) {

	if len(req.Include) == 0 {
		req.Include = []Include{
			IncludeDistances,
			IncludeMetadatas,
			IncludeDocuments,
		}
	}

	res := QueryResult{}
	err := v.send(
		ctx,
		fmt.Sprintf("/api/v1/collections/%s/query", collectionID),
		http.MethodPost,
		&req,
		&res,
	)
	return res, err
}

func (v *Client) send(ctx context.Context, path string, method string, obj any, out any) error {

	buffer := bytes.NewBuffer(nil)

	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("err: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf(
			"%s/%s",
			strings.TrimSuffix(v.url, "/"),
			strings.TrimPrefix(path, "/"),
		),
		buffer,
	)
	if err != nil {
		return fmt.Errorf("unable to prepare request: %w", err)
	}

	resp, err := v.client.Do(request)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		content, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server was unable to process %s: %s\n\n%s", path, resp.Status, content)
	}

	if out == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
