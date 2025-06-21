package thread

import (
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client/thread"
)

/*** Get thread response ***/

type getThreadResponseRequest struct {
	ClientID string                   `json:"client_id" validate:"required"`
	Thread   thread.Thread            `json:"thread"`
	Params   adapter.CompletionParams `json:"params"`
}

type getThreadResponseResponse struct {
	Thread thread.Thread `json:"thread"`
}

/*** Get thread completion ***/

type getThreadCompletionRequest struct {
	ClientID string                   `json:"client_id" validate:"required"`
	Thread   thread.Thread            `json:"thread"`
	Params   adapter.CompletionParams `json:"params"`
}

type getThreadCompletionResponse struct {
	thread.Completion
}

/*** Get thread summary ***/

type getThreadSummaryRequest struct {
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

type getThreadSummaryResponse struct {
	Thread thread.Thread `json:"thread"`
}
