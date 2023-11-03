package embeddings

// EmbeddingConfig represents the embedding config.
type EmbeddingConfig struct {
	Model     string
	BatchSize int
	Debug     bool
}

// Option define the function to set embeddings options.
type Option func(*EmbeddingConfig)

// OptionBathSize set the batch size (default to embbeding enginne default).
func OptionBathSize(b int) Option {
	return func(c *EmbeddingConfig) {
		c.BatchSize = b
	}
}

// OptionDebug set the debug mode (default: false).
func OptionDebug(debug bool) Option {
	return func(c *EmbeddingConfig) {
		c.Debug = debug
	}
}

// OptionModel set the model to use.
func OptionModel(model string) Option {
	return func(c *EmbeddingConfig) {
		c.Model = model
	}
}
