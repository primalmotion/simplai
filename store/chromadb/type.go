package chromadb

// Metadata holds chromadb metadata.
type Metadata map[string]any

// Where holds chromadb where clause.
type Where map[string]string

// Embedding holds chromadb embeddings.
// This only supports float64 for now.
type Embedding []float64

// Include represents the availaible include
// for Get and Query.
type Include string

const (
	IncludeMetadatas = "metadata"
	IncludeDocuments = "documents"
	IncludeDistances = "distances"
)

// EmbeddingAdd represents an embedding add request.
type EmbeddingAdd struct {
	Embeddings []Embedding `json:"embeddings,omitempty"`
	Metadatas  []Metadata  `json:"metadatas,omitempty"`
	Documents  []string    `json:"documents,omitempty"`
	IDs        []string    `json:"ids,omitempty"`
}

// EmbeddingUpdate represents an embedding update request.
type EmbeddingUpdate struct {
	Embeddings []Embedding `json:"embeddings,omitempty"`
	Metadatas  []Metadata  `json:"metadatas,omitempty"`
	Documents  []string    `json:"documents,omitempty"`
	IDs        []string    `json:"ids,omitempty"`
}

// EmbeddingQuery represents an embedding query request.
type EmbeddingQuery struct {
	Where           Where       `json:"where,omitempty"`
	WhereDocument   Where       `json:"where_documents,omitempty"`
	QueryEmbeddings []Embedding `json:"query_embeddings,omitempty"`
	Include         []Include   `json:"include,omitempty"`
	NResults        int         `json:"n_results,omitempty"`
}

// EmbeddingGet represents an embedding get request.
type EmbeddingGet struct {
	Where         Where     `json:"where,omitempty"`
	WhereDocument Where     `json:"where_documents,omitempty"`
	Sort          string    `json:"sort,omitempty"`
	IDs           []string  `json:"ids,omitempty"`
	Include       []Include `json:"include,omitempty"`
	Limit         int       `json:"limit,omitempty"`
	Offset        int       `json:"offset,omitempty"`
}

// EmbeddingDelete represents an embedding delete request.
type EmbeddingDelete struct {
	Where         Where    `json:"where,omitempty"`
	WhereDocument Where    `json:"where_documents,omitempty"`
	IDs           []string `json:"ids,omitempty"`
}

// EmbeddingCreate represents a collection create request.
type CollectionCreate struct {
	Metadata    Metadata `json:"metadata,omitempty"`
	Name        string   `json:"name"`
	GetOrCreate bool     `json:"get_or_create,omitempty"`
}

// CollectionUpdate represents a collection update request.
type CollectionUpdate struct {
	NewMetadata Metadata `json:"new_metadata,omitempty"`
	NewName     string   `json:"new_name"`
}

// GetResult holds the data returned by an EmbeddingGet request.
type GetResult struct {
	IDs        []string    `json:"ids,omitempty"`
	Embeddings []Embedding `json:"embeddings,omitempty"`
	Metadatas  []Metadata  `json:"metadatas,omitempty"`
	Documents  []string    `json:"documents,omitempty"`
}

// GetResult holds the data returned by an EmbeddingQuery request.
type QueryResult struct {
	IDs        [][]string  `json:"ids,omitempty"`
	Embeddings []Embedding `json:"embeddings,omitempty"`
	Metadatas  []Metadata  `json:"metadatas,omitempty"`
	Documents  [][]string  `json:"documents,omitempty"`
	Distances  [][]float64 `json:"distances,omitempty"`
}

// CollectionResult holds the data returned by an CollectionCreate request.
type CollectionResult struct {
	Metadatas Metadata `json:"metadatas,omitempty"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
}
