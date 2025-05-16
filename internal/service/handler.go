package service

import (
	"github.com/umk/llmservices/internal/jsonrpc"
	"github.com/umk/llmservices/internal/service/handlers/client"
	"github.com/umk/llmservices/internal/service/handlers/vectors"
)

func Handler() *jsonrpc.Handler {
	return jsonrpc.NewHandler(map[string]jsonrpc.HandlerFunc{
		"createDatabase": vectors.CreateDatabase,
		"deleteDatabase": vectors.DeleteDatabase,

		"addVector":          vectors.AddVector,
		"deleteVector":       vectors.DeleteVector,
		"addVectorsBatch":    vectors.AddVectorsBatch,
		"deleteVectorsBatch": vectors.DeleteVectorsBatch,
		"searchVectors":      vectors.SearchVectors,

		"getSimilarity": vectors.GetSimilarity,

		"setClient":           client.SetClient,
		"getCompletion":       client.GetCompletion,
		"getThreadCompletion": client.GetThreadCompletion,
		"getEmbeddings":       client.GetEmbeddings,
		"getStatistics":       client.GetStatistics,
	})
}
