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

	resp, err := cl.Completion(ctx, req.Messages, req.Params)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.GetResponse(getCompletionResponse{
		Completion: resp,
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

	resp, err := cl.Embeddings(ctx, req.Input, req.Params)
	if err != nil {
		return nil, newEmbeddingsError(err)
	}

	return c.GetResponse(getEmbeddingsResponse{
		Embeddings: resp,
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

func GetThreadCompletion(ctx context.Context, c jsonrpc.RPCContext) (any, error) {
	var req getThreadCompletionRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientId)
	if cl == nil {
		return nil, errClientNotFound
	}

	resp, err := cl.ThreadCompletion(ctx, &req.Thread, req.Params)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.GetResponse(getThreadCompletionResponse{
		ThreadCompletion: resp,
	})
}

func getClient(clientId string) *client.Client {
	v, ok := clients.Load(clientId)
	if !ok {
		return nil
	}

	return v.(*client.Client)
}
