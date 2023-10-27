package llm

type InferenceConfig struct {
	Model       string
	MaxTokens   int
	Temperature float64
	TopK        float64
}

func NewInferenceConfig() InferenceConfig {
	return InferenceConfig{
		MaxTokens:   512,
		Temperature: 0.0,
		TopK:        0.1,
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

func OptionInferTopK(topK float64) InferenceOption {
	return func(c *InferenceConfig) {
		c.TopK = topK
	}
}

func OptionInferMaxTokens(maxTokens int) InferenceOption {
	return func(c *InferenceConfig) {
		c.MaxTokens = maxTokens
	}
}
