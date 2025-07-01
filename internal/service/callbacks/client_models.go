package callbacks

import "github.com/umk/llmservices/pkg/adapter"

/*** Get function call ***/

type GetFunctionCallRequest struct {
	adapter.ToolCallFunction
}

type GetFunctionCallResponse struct {
	Response string `json:"response"`
}

/*** Push thought ***/

type PushThoughtRequest struct {
	Content string `json:"content"`
}
