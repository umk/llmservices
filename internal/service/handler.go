package service

import (
	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/service/handlers"
	"github.com/umk/llmservices/internal/service/handlers/agent"
	"github.com/umk/llmservices/internal/service/handlers/thread"
)

func Handler() *jsonrpc2.Handler {
	return jsonrpc2.NewHandler(map[string]jsonrpc2.HandlerFunc{
		"setClient":     handlers.SetClientRPC,
		"getCompletion": handlers.GetCompletionRPC,
		"getEmbeddings": handlers.GetEmbeddingsRPC,
		"getSpeech":     handlers.GetSpeechRPC,
		"getStatistics": handlers.GetStatisticsRPC,

		"getThreadCompletion": thread.GetCompletionRPC,
		"getThreadSummary":    thread.GetSummaryRPC,
		"getThreadResponse":   thread.GetResponseRPC,

		"getAgentResponse": agent.GetResponseRPC,
	})
}
