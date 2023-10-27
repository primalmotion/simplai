package llm

type InferenceConfig struct {
	LogitBias        map[string]int
	Model            string
	Stop             []string
	MaxTokens        int
	Temperature      float64
	FrequencyPenalty float64
	PresencePenalty  float64
	LogProbs         int
	TopP             float64
}

func NewInferenceConfig() InferenceConfig {
	return InferenceConfig{
		MaxTokens:   512,
		Temperature: 0.0,
		TopP:        1,
	}
}

type Option func(*InferenceConfig)

func OptionTemperature(temp float64) Option {
	return func(c *InferenceConfig) {
		c.Temperature = temp
	}
}

func OptionModel(model string) Option {
	return func(c *InferenceConfig) {
		c.Model = model
	}
}

func OptionMaxTokens(maxTokens int) Option {
	return func(c *InferenceConfig) {
		c.MaxTokens = maxTokens
	}
}

func OptionFrequencePenalty(penalty float64) Option {
	return func(c *InferenceConfig) {
		c.FrequencyPenalty = penalty
	}
}

func OptionLogitBias(bias map[string]int) Option {
	return func(c *InferenceConfig) {
		c.LogitBias = bias
	}
}

func OptionPresencePenalty(penalty float64) Option {
	return func(c *InferenceConfig) {
		c.PresencePenalty = penalty
	}
}

func OptionLogProbs(prob int) Option {
	return func(c *InferenceConfig) {
		c.LogProbs = prob
	}
}

func OptionStop(words ...string) Option {
	return func(c *InferenceConfig) {
		c.Stop = words
	}
}

func OptionTopP(topP float64) Option {
	return func(c *InferenceConfig) {
		c.TopP = topP
	}
}
