package service

import (
	"github.com/umk/llmservices/internal/jsonrpc"
)

func Handler() *jsonrpc.Handler {
	return jsonrpc.NewHandler(map[string]jsonrpc.HandlerFunc{
		"createDatabase": createDatabase,
		"deleteDatabase": deleteDatabase,

		"addVector":          addVector,
		"deleteVector":       deleteVector,
		"addVectorsBatch":    addVectorsBatch,
		"deleteVectorsBatch": deleteVectorsBatch,
		"searchVectors":      searchVectors,

		"getSimilarity": getSimilarity,
	})
}
