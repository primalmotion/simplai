package llm

// EmbeddingConfig represents the embedding config.
type EmbeddingConfig struct {
	Model     string
	BatchSize int
	Debug     bool
}

// EmbeddingOption define the function to set embeddings options.
type EmbeddingOption func(*EmbeddingConfig)

// OptionEmbeddingBatchSize set the batch size (default to embbeding enginne default).
func OptionEmbeddingBatchSize(b int) EmbeddingOption {
	return func(c *EmbeddingConfig) {
		c.BatchSize = b
	}
}

// OptionEmbeddingModel set the model to use.
func OptionEmbeddingModel(model string) EmbeddingOption {
	return func(c *EmbeddingConfig) {
		c.Model = model
	}
}

// OptionEmbeddingDebug set the debug mode (default: false).
func OptionEmbeddingDebug(debug bool) EmbeddingOption {
	return func(c *EmbeddingConfig) {
		c.Debug = debug
	}
}
