package chunker

// Options are the Options for a chunker.
type Options struct {
	Separators          []string
	ChunkSize           int
	ChunkOverlapPercent int
	TokenRatio          float32
}

// DefaultOptions set DefaultOptions.
func DefaultOptions() Options {
	return Options{
		Separators:          []string{"\n\n", "\n", ". ", " ", ""},
		ChunkSize:           512,
		ChunkOverlapPercent: 15,
		TokenRatio:          1.0,
	}
}

// Option is a function that can be used to set options for a chunker.
type Option func(*Options)

// OptionsChunkSize sets the chunk size for a chunker (default: 512 tokens).
func OptionsChunkSize(c int) Option {
	return func(o *Options) {
		o.ChunkSize = c
	}
}

// OptionsChunkOverlapPercent sets the chunk overlap percentage for a chunker (default: 15%).
func OptionsChunkOverlap(pct int) Option {
	return func(o *Options) {
		o.ChunkOverlapPercent = pct
	}
}

// OptionsSeparators sets the separators for a chunker.
func OptionsSeparators(s []string) Option {
	return func(o *Options) {
		o.Separators = s
	}
}

// OptionsTokenRatio sets the token ratio for a chunker.
// Ideal the token count should be done via the model tokeniser,
// in general for llm its around 100 token for 75 word, for BERT models
// its more around 1 (default: 1).
func OptionsTokenRatio(r float32) Option {
	return func(o *Options) {
		o.TokenRatio = r
	}
}
