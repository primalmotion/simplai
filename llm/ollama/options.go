package ollama

import (
	"log"
	"net/url"

	ollamaclient "git.sr.ht/~primalmotion/simplai/llm/ollama/internal"
)

type options struct {
	ollamaServerURL     *url.URL
	model               string
	customModelTemplate string
	system              string
	ollamaOptions       ollamaclient.Options
}

type Option func(*options)

// OptionModel Set the model to use.
func OptionModel(model string) Option {
	return func(opts *options) {
		opts.model = model
	}
}

// OptionSystem Set the system prompt. This is only valid if
// OptionCustomTemplate is not set and the ollama model use
// .System in its model template OR if OptionCustomTemplate
// is set using {{.System}}.
func OptionSystemPrompt(p string) Option {
	return func(opts *options) {
		opts.system = p
	}
}

// OptionCustomTemplate To override the templating done on Ollama model side.
func OptionCustomTemplate(template string) Option {
	return func(opts *options) {
		opts.customModelTemplate = template
	}
}

// OptionServerURL Set the URL of the ollama instance to use.
func OptionServerURL(rawURL string) Option {
	return func(opts *options) {
		var err error
		opts.ollamaServerURL, err = url.Parse(rawURL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// OptionRunnerNumKeep Specify the number of tokens from the initial prompt to retain when the model resets
// its internal context.
func OptionRunnerNumKeep(num int) Option {
	return func(opts *options) {
		opts.ollamaOptions.NumKeep = num
	}
}

// OptionRunnerNumThread Set the number of threads to use during computation (default: auto).
func OptionRunnerNumThread(num int) Option {
	return func(opts *options) {
		opts.ollamaOptions.NumThread = num
	}
}

// OptionRunnerNumGPU The number of layers to send to the GPU(s).
// On macOS it defaults to 1 to enable metal support, 0 to disable.
func OptionRunnerNumGPU(num int) Option {
	return func(opts *options) {
		opts.ollamaOptions.NumGPU = num
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
