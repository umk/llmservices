package vectors

import "github.com/umk/jsonrpc2"

var errDatabaseNotFound = jsonrpc2.Error{
	Code:    -32000,
	Message: "Database not found",
}

var errDatabaseAlreadyExists = jsonrpc2.Error{
	Code:    -32000,
	Message: "Database already exists",
}

var errVectorsLengthMismatch = jsonrpc2.Error{
	Code:    -32000,
	Message: "Vectors must have the same length",
}
