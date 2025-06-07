package client

import "github.com/umk/llmservices/pkg/client"

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
