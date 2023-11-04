package engine

// InferenceConfig represents the inference configuration for a LLM.
type InferenceConfig struct {
	LogitBias         map[string]int
	Model             string
	Stop              []string
	MaxTokens         int
	Temperature       float64
	FrequencyPenalty  float64
	RepetitionPenalty float64
	PresencePenalty   float64
	LogProbs          int
	TopP              float64
	TopK              int
	Seed              int
	Debug             bool
}

// LLMOption define the function to set InferenceConfig options.
type LLMOption func(*InferenceConfig)

// OptionDebug set the debug mode (default: false).
func OptionDebug(debug bool) LLMOption {
	return func(c *InferenceConfig) {
		c.Debug = debug
	}
}

// OptionTemperature set the Temperature.
func OptionTemperature(temp float64) LLMOption {
	return func(c *InferenceConfig) {
		c.Temperature = temp
	}
}

// OptionModel set the model to use.
func OptionModel(model string) LLMOption {
	return func(c *InferenceConfig) {
		c.Model = model
	}
}

// OptionMaxTokens set the maximun number of token to infer.
func OptionMaxTokens(maxTokens int) LLMOption {
	return func(c *InferenceConfig) {
		c.MaxTokens = maxTokens
	}
}

// OptionLogitBias Modify the likelihood of a token appearing in the generated text completion.
func OptionLogitBias(bias map[string]int) LLMOption {
	return func(c *InferenceConfig) {
		c.LogitBias = bias
	}
}

// OptionRepetitionPenalty set the repetition penalty.
func OptionRepetitionPenalty(penalty float64) LLMOption {
	return func(c *InferenceConfig) {
		c.RepetitionPenalty = penalty
	}
}

// OptionPresencePenalty set the presence penalty
func OptionPresencePenalty(penalty float64) LLMOption {
	return func(c *InferenceConfig) {
		c.PresencePenalty = penalty
	}
}

// OptionFrequencePenalty set the frequencies penalty.
func OptionFrequencePenalty(penalty float64) LLMOption {
	return func(c *InferenceConfig) {
		c.FrequencyPenalty = penalty
	}
}

// OptionLogProbs returns the logarithm of the density or probability.
func OptionLogProbs(prob int) LLMOption {
	return func(c *InferenceConfig) {
		c.LogProbs = prob
	}
}

// OptionStop sets the stop words.
func OptionStop(words ...string) LLMOption {
	return func(c *InferenceConfig) {
		c.Stop = words
	}
}

// OptionTopP limits the next token selection to a subset of tokens with
// a cumulative probability above a threshold P.
func OptionTopP(topP float64) LLMOption {
	return func(c *InferenceConfig) {
		c.TopP = topP
	}
}

// OptionTopK limits the next token selection to the K most probable tokens.
func OptionTopK(topK int) LLMOption {
	return func(c *InferenceConfig) {
		c.TopK = topK
	}
}

// OptionSeed  sets the random number generator (RNG) seed (default: -1, -1 = random seed).
func OptionSeed(seed int) LLMOption {
	return func(c *InferenceConfig) {
		c.Seed = seed
	}
}
