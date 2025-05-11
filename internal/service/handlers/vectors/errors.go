package vectors

import "github.com/umk/llmservices/internal/jsonrpc"

var errDatabaseNotFound = jsonrpc.Error{
	Code:    -32000,
	Message: "Database not found",
}

var errDatabaseAlreadyExists = jsonrpc.Error{
	Code:    -32000,
	Message: "Database already exists",
}

var errVectorsLengthMismatch = jsonrpc.Error{
	Code:    -32000,
	Message: "Vectors must have the same length",
}
