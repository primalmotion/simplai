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

type InferenceOption func(*InferenceConfig)

func OptionInferTemperature(temp float64) InferenceOption {
	return func(c *InferenceConfig) {
		c.Temperature = temp
	}
}

func OptionInferModel(model string) InferenceOption {
	return func(c *InferenceConfig) {
		c.Model = model
	}
}

func OptionInferMaxTokens(maxTokens int) InferenceOption {
	return func(c *InferenceConfig) {
		c.MaxTokens = maxTokens
	}
}

func OptionInferFrequencePenalty(penalty float64) InferenceOption {
	return func(c *InferenceConfig) {
		c.FrequencyPenalty = penalty
	}
}

func OptionInferLogitBias(bias map[string]int) InferenceOption {
	return func(c *InferenceConfig) {
		c.LogitBias = bias
	}
}

func OptionInferPresencePenalty(penalty float64) InferenceOption {
	return func(c *InferenceConfig) {
		c.PresencePenalty = penalty
	}
}

func OptionInferLogProbs(prob int) InferenceOption {
	return func(c *InferenceConfig) {
		c.LogProbs = prob
	}
}

func OptionInferStop(stop ...string) InferenceOption {
	return func(c *InferenceConfig) {
		c.Stop = stop
	}
}

func OptionInferTopP(topP float64) InferenceOption {
	return func(c *InferenceConfig) {
		c.TopP = topP
	}
}
