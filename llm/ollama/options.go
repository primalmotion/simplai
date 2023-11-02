package ollama

import (
	"github.com/primalmotion/simplai/llm"
	ollamaclient "github.com/primalmotion/simplai/llm/ollama/internal"
)

type options struct {
	customModelTemplate    string
	system                 string
	defaultInferenceConfig llm.InferenceConfig
	ollamaOptions          ollamaclient.Options
	raw                    bool
}

// Option is the function to handle options.
type Option func(*options)

func defaultOptions() options {
	return options{
		ollamaOptions: ollamaclient.DefaultOptions(),
		raw:           true,
	}
}

// OptionDefaultInferenceConfig To set the default InferenceConfig parameters.
func OptionDefaultInferenceConfig(c llm.InferenceConfig) Option {
	return func(opts *options) {
		opts.defaultInferenceConfig = c
	}
}

// OptionSystemPrompt Set the system prompt. This is only valid if
// OptionCustomTemplate is not set and the ollama model use
// .System in its model template OR if OptionCustomTemplate
// is set using {{.System}}.
func OptionSystemPrompt(p string) Option {
	return func(opts *options) {
		opts.system = p
	}
}

// OptionUseModelTemplating To enable the prompt templating done
// on the ollama side if set in the modelfile (default is false).
func OptionUseModelTemplating(b bool) Option {
	return func(opts *options) {
		opts.raw = !b
	}
}

// OptionCustomTemplate To override the templating done on Ollama model side.
func OptionCustomTemplate(template string) Option {
	return func(opts *options) {
		opts.customModelTemplate = template
	}
}

// OptionRunnerNumKeep Specify the number of tokens from the initial prompt to retain when the model resets
// its internal context.
func OptionRunnerNumKeep(num int) Option {
	return func(opts *options) {
		opts.ollamaOptions.NumKeep = num
	}
}

// OptionPredictTFSZ Tail free sampling is used to reduce the impact of less probable tokens from the output.
// A higher value (e.g., 2.0) will reduce the impact more, while a value of 1.0 disables this setting (default: 1).
func OptionPredictTFSZ(val float32) Option {
	return func(opts *options) {
		opts.ollamaOptions.TFSZ = val
	}
}

// OptionPredictTypicalP Enable locally typical sampling with parameter p (default: 1.0, 1.0 = disabled).
func OptionPredictTypicalP(val float32) Option {
	return func(opts *options) {
		opts.ollamaOptions.TypicalP = val
	}
}

// OptionPredictRepeatLastN Sets how far back for the model to look back to prevent repetition
// (Default: 64, 0 = disabled, -1 = num_ctx).
func OptionPredictRepeatLastN(val int) Option {
	return func(opts *options) {
		opts.ollamaOptions.RepeatLastN = val
	}
}

// OptionPredictMirostat Enable Mirostat sampling for controlling perplexity
// (default: 0, 0 = disabled, 1 = Mirostat, 2 = Mirostat 2.0).
func OptionPredictMirostat(val int) Option {
	return func(opts *options) {
		opts.ollamaOptions.Mirostat = val
	}
}

// OptionPredictMirostatTau Controls the balance between coherence and diversity of the output.
// A lower value will result in more focused and coherent text (Default: 5.0).
func OptionPredictMirostatTau(val float32) Option {
	return func(opts *options) {
		opts.ollamaOptions.MirostatTau = val
	}
}

// OptionPredictMirostatEta Influences how quickly the algorithm responds to feedback from the generated text.
// A lower learning rate will result in slower adjustments, while a higher learning rate will make the
// algorithm more responsive (Default: 0.1).
func OptionPredictMirostatEta(val float32) Option {
	return func(opts *options) {
		opts.ollamaOptions.MirostatEta = val
	}
}

// OptionPredictPenalizeNewline Penalize newline tokens when applying the repeat penalty (default: true).
func OptionPredictPenalizeNewline(val bool) Option {
	return func(opts *options) {
		opts.ollamaOptions.PenalizeNewline = val
	}
}
