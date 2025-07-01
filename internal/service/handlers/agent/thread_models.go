package agent

import (
	"github.com/umk/llmservices/pkg/client/agent"
	"github.com/umk/llmservices/pkg/client/thread"
)

/*** Get response ***/

type GetResponseRequest struct {
	ClientID string               `json:"client_id" validate:"required"`
	Thread   thread.Thread        `json:"thread"`
	Params   agent.ResponseParams `json:"params"`
}

type GetResponseResponse struct {
	agent.Response
}
