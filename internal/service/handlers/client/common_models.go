package client

import (
	"github.com/umk/llmservices/pkg/adapter"
)

/*** Get completion ***/

type getCompletionRequest struct {
	ClientID string                   `json:"client_id" validate:"required"`
	Messages []adapter.Message        `json:"messages" validate:"required,min=1"`
	Params   adapter.CompletionParams `json:"params"`
}

type getCompletionResponse struct {
	adapter.Completion
}

/*** Get embeddings ***/

type getEmbeddingsRequest struct {
	Input    string                   `json:"input" validate:"required"`
	Params   adapter.EmbeddingsParams `json:"params"`
	ClientID string                   `json:"client_id" validate:"required"`
}

type getEmbeddingsResponse struct {
	adapter.Embeddings
}

/*** Get statistics ***/

type getStatisticsRequest struct {
	ClientID string `json:"client_id" validate:"required"`
}

type getStatisticsResponse struct {
	BytesPerTok float32 `json:"bytes_per_tok"`
}
