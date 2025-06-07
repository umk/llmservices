package service

import (
	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/service/handlers/client"
	"github.com/umk/llmservices/internal/service/handlers/vectors"
)

func Handler() *jsonrpc2.Handler {
	return jsonrpc2.NewHandler(map[string]jsonrpc2.HandlerFunc{
		"createDatabase": vectors.CreateDatabase,
		"deleteDatabase": vectors.DeleteDatabase,
		"readDatabase":   vectors.ReadDatabase,
		"writeDatabase":  vectors.WriteDatabase,

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
