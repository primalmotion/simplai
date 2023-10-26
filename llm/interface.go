package llm

// LLMN is the main interface to interfact with a LLM.
type LLM interface {
	Infer(text string, options ...InferenceOption) (string, error)
}

type InferenceConfig struct {
	Model       string
	MaxTokens   int
	Temperature float32
	TopK        float32
}

func NewInferenceConfig() InferenceConfig {
	return InferenceConfig{
		MaxTokens:   512,
		Temperature: 0.0,
		TopK:        0.1,
	}
}

type InferenceOption func(*InferenceConfig)

func OptionInferTemperature(temp float32) InferenceOption {
	return func(c *InferenceConfig) {
		c.Temperature = temp
	}
}

func OptionInferModel(model string) InferenceOption {
	return func(c *InferenceConfig) {
		c.Model = model
	}
}

func OptionInferTopK(topK float32) InferenceOption {
	return func(c *InferenceConfig) {
		c.TopK = topK
	}
}

func OptionInferMaxTokens(maxTokens int) InferenceOption {
	return func(c *InferenceConfig) {
		c.MaxTokens = maxTokens
	}
}
