package client

import "github.com/umk/llmservices/pkg/adapter"

/*** Generate Speech ***/

type getSpeechRequest struct {
	ClientId string                `json:"client_id" validate:"required"`
	Message  adapter.SpeechMessage `json:"message"`
	Params   adapter.SpeechParams  `json:"params"`
}

type getSpeechResponse struct {
	adapter.Speech
}
