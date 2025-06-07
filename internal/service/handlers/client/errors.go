package client

import "github.com/umk/jsonrpc2"

var errClientNotFound = jsonrpc2.Error{
	Code:    -32000,
	Message: "Client not found",
}

func newCompletionError(err error) error {
	return jsonrpc2.Error{
		Code:    -32000,
		Message: "Completion error",
		Data:    map[string]any{"error": err.Error()},
	}
}

func newEmbeddingsError(err error) error {
	return jsonrpc2.Error{
		Code:    -32000,
		Message: "Embeddings error",
		Data:    map[string]any{"error": err.Error()},
	}
}
