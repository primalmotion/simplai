package openai

import (
	"github.com/primalmotion/simplai/llm"
)

type options struct {
	defaultInferenceConfig llm.InferenceConfig
}

// Option represent the func option.
type Option func(*options)

func defaultOptions() options {
	return options{
		defaultInferenceConfig: llm.InferenceConfig{
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

func defaultEmbeddingConfig() llm.EmbeddingConfig {
	return llm.EmbeddingConfig{
		BatchSize: 512,
	}
}

// OptionDefaultInferenceConfig To set the default InferenceConfig parameters.
func OptionDefaultInferenceConfig(c llm.InferenceConfig) Option {
	return func(opts *options) {
		opts.defaultInferenceConfig = c
	}
}
