package api

import (
	"context"
	"encoding/json"
)

type ClientType int

const (
	Gemini_2_5 ClientType = iota
	OllamaQwen
)

type CompletionRequest struct {
	method string 
	endpoint string
	Body any
}

type CompletionResponse struct {
	Body string
}

// TODO: have google, ollama, both have clients that implement this function
// Use Req and Resp generics to support multiple implementations of CompletionClient
type CompletionClient[Req any, Resp any] interface {
	GenerateCompletion(ctx context.Context, req Req, fn func(Resp) error) error
}	

/// Ollama types
// GenerateRequest describes a request sent by [Client.Generate]. While you
// have to specify the Model and Prompt fields, all the other fields have
// reasonable defaults for basic uses.
type GenerateRequest struct {
	// Model is the model name; it should be a name familiar to Ollama from
	// the library at https://ollama.com/library
	Model string `json:"model"`

	// Prompt is the textual prompt to send to the model.
	Prompt string `json:"prompt"`

	// Suffix is the text that comes after the inserted text.
	Suffix string `json:"suffix"`

	// System overrides the model's default system message/prompt.
	System string `json:"system"`

	// Template overrides the model's default prompt template.
	Template string `json:"template"`

	// Context is the context parameter returned from a previous call to
	// [Client.Generate]. It can be used to keep a short conversational memory.
	Context []int `json:"context,omitempty"`

	// Stream specifies whether the response is streaming; it is true by default.
	Stream *bool `json:"stream,omitempty"`

	// Raw set to true means that no formatting will be applied to the prompt.
	Raw bool `json:"raw,omitempty"`

	// Format specifies the format to return a response in.
	Format json.RawMessage `json:"format,omitempty"`

	// KeepAlive controls how long the model will stay loaded in memory following
	// this request.
	// KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Images is an optional list of raw image bytes accompanying this
	// request, for multimodal models.
	// Images []ImageData `json:"images,omitempty"`

	// Options lists model-specific options. For example, temperature can be
	// set through this field, if the model supports it.
	Options map[string]any `json:"options"`

	// Think controls whether thinking/reasoning models will think before
	// responding. Can be a boolean (true/false) or a string ("high", "medium", "low")
	// for supported models. Needs to be a pointer so we can distinguish between false
	// (request that thinking _not_ be used) and unset (use the old behavior
	// before this option was introduced)
	// Think *ThinkValue `json:"think,omitempty"`

	// Truncate is a boolean that, when set to true, truncates the chat history messages
	// if the rendered prompt exceeds the context length limit.
	Truncate *bool `json:"truncate,omitempty"`

	// Shift is a boolean that, when set to true, shifts the chat history
	// when hitting the context length limit instead of erroring.
	Shift *bool `json:"shift,omitempty"`

	// DebugRenderOnly is a debug option that, when set to true, returns the rendered
	// template instead of calling the model.
	DebugRenderOnly bool `json:"_debug_render_only,omitempty"`

	// Logprobs specifies whether to return log probabilities of the output tokens.
	Logprobs bool `json:"logprobs,omitempty"`

	// TopLogprobs is the number of most likely tokens to return at each token position,
	// each with an associated log probability. Only applies when Logprobs is true.
	// Valid values are 0-20. Default is 0 (only return the selected token's logprob).
	TopLogprobs int `json:"top_logprobs,omitempty"`

	// Experimental: Image generation fields (may change or be removed)

	// Width is the width of the generated image in pixels.
	// Only used for image generation models.
	Width int32 `json:"width,omitempty"`

	// Height is the height of the generated image in pixels.
	// Only used for image generation models.
	Height int32 `json:"height,omitempty"`

	// Steps is the number of diffusion steps for image generation.
	// Only used for image generation models.
	Steps int32 `json:"steps,omitempty"`
}

type GenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Thinking           string `json:"thinking"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	TotalDuration      int    `json:"total_duration"`
	LoadDuration       int    `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int    `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int    `json:"eval_duration"`
}