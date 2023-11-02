package openai

import "git.sr.ht/~primalmotion/simplai/llm"

type options struct {
	defaultInferenceConfig llm.InferenceConfig
}

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

// OptionDefaultInferenceConfig To set the default InferenceConfig parameters.
func OptionDefaultInferenceConfig(c llm.InferenceConfig) Option {
	return func(opts *options) {
		opts.defaultInferenceConfig = c
	}
}
