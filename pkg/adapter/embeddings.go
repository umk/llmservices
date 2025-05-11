package adapter

type EmbeddingsRequest struct {
	Input      string `json:"input" validate:"required"`
	Model      string `json:"model" validate:"required"`
	Dimensions *int64 `json:"dimensions,omitempty" validate:"omitempty,min=1"`
}

type EmbeddingsResponse struct {
	Data  []float64        `json:"data" validate:"required,min=1"`
	Usage *EmbeddingsUsage `json:"usage,omitempty"`
}

type EmbeddingsUsage struct {
	PromptTokens int64 `json:"prompt_tokens"`
}
