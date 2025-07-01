package agent

import "github.com/umk/jsonrpc2"

func newResponseError(err error) error {
	return jsonrpc2.Error{
		Code:    -32000,
		Message: "Response error",
		Data:    map[string]any{"error": err.Error()},
	}
}
