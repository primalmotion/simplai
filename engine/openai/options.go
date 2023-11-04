package openai

import (
	"github.com/primalmotion/simplai/engine"
)

type options struct {
	defaultInferenceConfig engine.InferenceConfig
}

// Option represent the func option.
type Option func(*options)

func defaultOptions() options {
	return options{
		defaultInferenceConfig: engine.InferenceConfig{
			Temperature:       1.0,
			FrequencyPenalty:  0,
			RepetitionPenalty: 1.0,
			PresencePenalty:   0,
			LogProbs:          0,
			TopP:              1,
			TopK:              0,
			Seed:              -1,
		},
	}
}

func defaultEmbeddingConfig() engine.EmbeddingConfig {
	return engine.EmbeddingConfig{
		BatchSize: 512,
	}
}

// OptionDefaultInferenceConfig To set the default InferenceConfig parameters.
func OptionDefaultInferenceConfig(c engine.InferenceConfig) Option {
	return func(opts *options) {
		opts.defaultInferenceConfig = c
	}
}
