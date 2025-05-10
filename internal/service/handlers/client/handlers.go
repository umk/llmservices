package client

import (
	"context"
	"sync"

	"github.com/umk/llmservices/internal/jsonrpc"
	"github.com/umk/llmservices/pkg/client"
)

var clients sync.Map

func SetClient(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req setClientRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	conf, err := getClientConfig(&req.Config)
	if err != nil {
		return nil, err
	}

	cl := client.New(conf)
	clients.Store(req.ClientId, cl)

	var resp setClientResponse

	return c.GetResponse(resp)
}

func GetCompletion(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req getCompletionRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientId)
	if cl == nil {
		return nil, errClientNotFound
	}

	resp, err := cl.GetCompletion(ctx, &req.CompletionRequest)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.GetResponse(getCompletionResponse{
		CompletionResponse: resp,
	})
}

func GetEmbeddings(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req getEmbeddingsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientId)
	if cl == nil {
		return nil, errClientNotFound
	}

	resp, err := cl.GetEmbeddings(ctx, &req.EmbeddingsRequest)
	if err != nil {
		return nil, newEmbeddingsError(err)
	}

	return c.GetResponse(getEmbeddingsResponse{
		EmbeddingsResponse: resp,
	})
}

func GetStatistics(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req getStatisticsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientId)
	if cl == nil {
		return nil, errClientNotFound
	}

	resp := getStatisticsResponse{
		BytesPerTok: cl.Samples.BytesPerTok(),
	}

	return c.GetResponse(resp)
}

func getClient(clientId string) *client.Client {
	v, ok := clients.Load(clientId)
	if !ok {
		return nil
	}

	return v.(*client.Client)
}
