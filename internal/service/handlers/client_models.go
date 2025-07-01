package handlers

import "github.com/umk/llmservices/pkg/client"

/*** Set client ***/

type SetClientRequest struct {
	ClientID string        `json:"client_id" validate:"required"`
	Config   client.Config `json:"config"`
}

type SetClientResponse struct{}
