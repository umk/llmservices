package client

import "github.com/umk/llmservices/pkg/adapter"

/*** Function call ***/

type functionCallRequest struct {
	adapter.ToolCallFunction
}

type functionCallResponse struct {
	Response string `json:"response"`
}
