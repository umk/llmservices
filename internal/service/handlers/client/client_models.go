package client

import "github.com/umk/llmservices/pkg/client"

/*** Set client ***/

type setClientRequest struct {
	ClientID string       `json:"client_id" validate:"required"`
	Config   clientConfig `json:"config"`
}

type clientConfig struct {
	Preset *client.Preset `json:"preset" validate:"omitempty,min=1"`

	BaseURL string `json:"base_url" validate:"omitempty,url"`
	Key     string `json:"key" validate:"omitempty"`
	Model   string `json:"model" validate:"omitempty"`

	Concurrency int `json:"concurrency" validate:"omitempty,min=1"`
}

type setClientResponse struct{}
