package handlers

import (
	"github.com/umk/llmservices/pkg/adapter"
)

/*** Get completion ***/

type GetCompletionRequest struct {
	ClientID string                   `json:"client_id" validate:"required"`
	Messages []adapter.Message        `json:"messages" validate:"required,min=1"`
	Params   adapter.CompletionParams `json:"params"`
}

type GetCompletionResponse struct {
	adapter.Completion
}

/*** Get embeddings ***/

type GetEmbeddingsRequest struct {
	Input    string                   `json:"input" validate:"required"`
	Params   adapter.EmbeddingsParams `json:"params"`
	ClientID string                   `json:"client_id" validate:"required"`
}

type GetEmbeddingsResponse struct {
	adapter.Embeddings
}

/*** Get statistics ***/

type GetStatisticsRequest struct {
	ClientID string `json:"client_id" validate:"required"`
}

type GetStatisticsResponse struct {
	BytesPerTok float32 `json:"bytes_per_tok"`
}
