package client

import (
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client"
)

/*** Set client ***/

type setClientRequest struct {
	ClientId string       `json:"client_id" validate:"required"`
	Config   clientConfig `json:"config"`
}

type clientConfig struct {
	Preset *client.Preset `yaml:"preset" validate:"omitempty,min=1"`

	BaseURL string `yaml:"base_url" validate:"omitempty,url"`
	Key     string `yaml:"key" validate:"omitempty"`
	Model   string `yaml:"model" validate:"omitempty"`

	Concurrency int `yaml:"concurrency" validate:"omitempty,min=1"`
}

type setClientResponse struct{}

/*** Get completion ***/

type getCompletionRequest struct {
	ClientId string                   `json:"client_id" validate:"required"`
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
	ClientId string                   `json:"client_id" validate:"required"`
}

type getEmbeddingsResponse struct {
	adapter.Embeddings
}

/*** Get statistics ***/

type getStatisticsRequest struct {
	ClientId string `json:"client_id" validate:"required"`
}

type getStatisticsResponse struct {
	BytesPerTok float32 `json:"bytes_per_tok"`
}

/*** Get thread completion ***/

type getThreadCompletionRequest struct {
	ClientId string                   `json:"client_id" validate:"required"`
	Thread   client.Thread            `json:"thread"`
	Params   adapter.CompletionParams `json:"params"`
}

type getThreadCompletionResponse struct {
	client.ThreadCompletion
}
