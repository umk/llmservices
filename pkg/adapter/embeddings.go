package adapter

type EmbeddingsParams struct {
	Model      string `json:"model" validate:"required"`
	Dimensions *int64 `json:"dimensions,omitempty" validate:"omitempty,min=1"`
}

type Embeddings struct {
	Data  []float64        `json:"data" validate:"required,min=1"`
	Usage *EmbeddingsUsage `json:"usage,omitempty"`
}

type EmbeddingsUsage struct {
	PromptTokens int64 `json:"prompt_tokens"`
}
