package service

import (
	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/service/handlers/client"
	"github.com/umk/llmservices/internal/service/handlers/client/thread"
	"github.com/umk/llmservices/internal/service/handlers/vectors"
)

func Handler() *jsonrpc2.Handler {
	return jsonrpc2.NewHandler(map[string]jsonrpc2.HandlerFunc{
		"createDatabase": vectors.CreateDatabaseRPC,
		"deleteDatabase": vectors.DeleteDatabaseRPC,
		"readDatabase":   vectors.ReadDatabaseRPC,
		"writeDatabase":  vectors.WriteDatabaseRPC,

		"addVector":          vectors.AddVectorRPC,
		"deleteVector":       vectors.DeleteVectorRPC,
		"addVectorsBatch":    vectors.AddVectorsBatchRPC,
		"deleteVectorsBatch": vectors.DeleteVectorsBatchRPC,
		"searchVectors":      vectors.SearchVectorsRPC,

		"getSimilarity": vectors.GetSimilarityRPC,

		"setClient":     client.SetClientRPC,
		"getCompletion": client.GetCompletionRPC,
		"getEmbeddings": client.GetEmbeddingsRPC,
		"getSpeech":     client.GetSpeechRPC,
		"getStatistics": client.GetStatisticsRPC,

		"getThreadCompletion": thread.GetCompletionRPC,
		"getThreadSummary":    thread.GetSummaryRPC,
		"getThreadResponse":   thread.GetResponseRPC,
	})
}
