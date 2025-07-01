package thread

import (
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client/thread"
)

/*** Get response ***/

type GetResponseRequest struct {
	ClientID string                `json:"client_id" validate:"required"`
	Thread   thread.Thread         `json:"thread"`
	Params   thread.ResponseParams `json:"params"`
}

type GetResponseResponse struct {
	thread.Response
}

/*** Get completion ***/

type GetCompletionRequest struct {
	ClientID string                   `json:"client_id" validate:"required"`
	Thread   thread.Thread            `json:"thread"`
	Params   adapter.CompletionParams `json:"params"`
}

type GetCompletionResponse struct {
	thread.Completion
}

/*** Get summary ***/

type GetSummaryRequest struct {
	// ID of a client that generates the summary
	ClientID string `json:"client_id" validate:"required"`
	// ID of a client that generated the thread completion. If not specified,
	// the generator is assumed to be the same as the summarizer.
	GenClientID *string       `json:"gen_client_id"`
	Thread      thread.Thread `json:"thread"`
	Fraction    float32       `json:"fraction"`

	MaxMessages *int `json:"max_messages" validate:"omitempty,gt=0"`
	MaxTokens   *int `json:"max_tokens" validate:"omitempty,gt=0"`
}

type GetSummaryResponse struct {
	Thread thread.Thread `json:"thread"`
}
