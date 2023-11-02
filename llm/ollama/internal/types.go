package ollamaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// StatusError describe various errors.
type StatusError struct {
	Status       string `json:"status,omitempty"`
	ErrorMessage string `json:"error"`
	StatusCode   int    `json:"code,omitempty"`
}

// Error return a proper error.
func (e StatusError) Error() string {
	switch {
	case e.Status != "" && e.ErrorMessage != "":
		return fmt.Sprintf("%s: %s", e.Status, e.ErrorMessage)
	case e.Status != "":
		return e.Status
	case e.ErrorMessage != "":
		return e.ErrorMessage
	default:
		// this should not happen
		return "something went wrong, please see the ollama server logs for details"
	}
}

// GenerateRequest represents an ollama request.
type GenerateRequest struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	System   string `json:"system"`
	Template string `json:"template"`
	Context  []int  `json:"context,omitempty"`
	Options  `json:"options"`
	Stream   bool `json:"stream"`
	Raw      bool `json:"raw"`
}

// String representation of an ollama request.
func (req GenerateRequest) String() string {
	req.Prompt = ""
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(req)
	return buf.String()
}

// GenerateResponse represents an ollama response.
type GenerateResponse struct {
	CreatedAt          time.Time     `json:"created_at"`
	Model              string        `json:"model"`
	Response           string        `json:"response"`
	Context            []int         `json:"context,omitempty"`
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
	Done               bool          `json:"done"`
}

// String representation of an ollama response.
func (resp GenerateResponse) String() string {
	resp.Context = []int{}
	resp.Response = ""
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	_ = encoder.Encode(resp)
	return buf.String()
}

// Options represents the ollama options for inference.
type Options struct {
	Stop             []string `json:"stop,omitempty"`
	RepeatLastN      int      `json:"repeat_last_n,omitempty"`
	Seed             int      `json:"seed"`
	TopK             int      `json:"top_k"`
	NumKeep          int      `json:"num_keep,omitempty"`
	Mirostat         int      `json:"mirostat"`
	NumPredict       int      `json:"num_predict,omitempty"`
	Temperature      float32  `json:"temperature"`
	TypicalP         float32  `json:"typical_p"`
	RepeatPenalty    float32  `json:"repeat_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	TFSZ             float32  `json:"tfs_z"`
	MirostatTau      float32  `json:"mirostat_tau"`
	MirostatEta      float32  `json:"mirostat_eta"`
	TopP             float32  `json:"top_p"`
	PenalizeNewline  bool     `json:"penalize_newline,omitempty"`
}

// DefaultOptions return the default ollama options.
func DefaultOptions() Options {
	return Options{
		NumPredict:       -1,
		NumKeep:          -1,
		Temperature:      0.8,
		TopK:             40,
		TopP:             0.9,
		TFSZ:             1.0,
		TypicalP:         1.0,
		RepeatLastN:      64,
		RepeatPenalty:    1.1,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		Mirostat:         0,
		MirostatTau:      5.0,
		MirostatEta:      0.1,
		PenalizeNewline:  true,
		Seed:             -1,
	}
}
