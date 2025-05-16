package adapter

type CompletionParams struct {
	Model            string          `json:"model" validate:"required"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty" validate:"omitempty,gte=-2.0,lte=2.0"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty" validate:"omitempty,gte=-2.0,lte=2.0"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Stop             []string        `json:"stop,omitempty" validate:"omitempty,max=4"`
	Temperature      *float64        `json:"temperature,omitempty" validate:"omitempty,gte=0.0,lte=2.0"`
	Tools            []Tool          `json:"tools,omitempty"`
	TopP             *float64        `json:"top_p,omitempty" validate:"omitempty,gte=0.0,lte=1.0"`
}

type Completion struct {
	Message AssistantMessage `json:"message" validate:"required"`
	Usage   *CompletionUsage `json:"usage,omitempty"`
}

type CompletionUsage struct {
	CompletionTokens int64 `json:"completion_tokens"`
	PromptTokens     int64 `json:"prompt_tokens"`
}

type ResponseFormat struct {
	OfResponseFormatText       *ResponseFormatText       `json:"text,omitempty"`
	OfResponseFormatJSONSchema *ResponseFormatJSONSchema `json:"json_schema,omitempty"`
}

type ResponseFormatText struct{}

type ResponseFormatJSONSchema struct {
	JSONSchema JSONSchema `json:"json_schema" validate:"required"`
}

type JSONSchema struct {
	Name        string         `json:"name" validate:"required"`
	Description *string        `json:"description,omitempty"`
	Schema      map[string]any `json:"schema" validate:"required"`
	Strict      *bool          `json:"strict,omitempty"`
}
