package thread

import "github.com/umk/jsonrpc2"

var errSummarizerParams = jsonrpc2.Error{
	Code:    -32000,
	Message: "Must specify either the max tokens or max messages",
}

func newResponseError(err error) error {
	return jsonrpc2.Error{
		Code:    -32000,
		Message: "Response error",
		Data:    map[string]any{"error": err.Error()},
	}
}

func newCompletionError(err error) error {
	return jsonrpc2.Error{
		Code:    -32000,
		Message: "Completion error",
		Data:    map[string]any{"error": err.Error()},
	}
}

func newSummarizerError(err error) error {
	return jsonrpc2.Error{
		Code:    -32000,
		Message: "Summarizer error",
		Data:    map[string]any{"error": err.Error()},
	}
}
