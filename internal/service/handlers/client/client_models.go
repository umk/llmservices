package client

import "github.com/umk/llmservices/pkg/client"

/*** Set client ***/

type setClientRequest struct {
	ClientID string        `json:"client_id" validate:"required"`
	Config   client.Config `json:"config"`
}

type setClientResponse struct{}
