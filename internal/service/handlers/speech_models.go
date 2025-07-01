package handlers

import "github.com/umk/llmservices/pkg/adapter"

/*** Generate Speech ***/

type GetSpeechRequest struct {
	ClientID string                `json:"client_id" validate:"required"`
	Message  adapter.SpeechMessage `json:"message"`
	Params   adapter.SpeechParams  `json:"params"`
}

type GetSpeechResponse struct {
	adapter.Speech
}
