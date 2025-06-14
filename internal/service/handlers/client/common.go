package client

import (
	"context"

	"github.com/umk/jsonrpc2"
)

func GetCompletion(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getCompletionRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientID)
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

func GetEmbeddings(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getEmbeddingsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientID)
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

func GetStatistics(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getStatisticsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientID)
	if cl == nil {
		return nil, errClientNotFound
	}

	resp := getStatisticsResponse{
		BytesPerTok: cl.Samples.BytesPerTok(),
	}

	return c.GetResponse(resp)
}

func GetThreadCompletion(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getThreadCompletionRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientID)
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
