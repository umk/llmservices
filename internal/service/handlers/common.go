package handlers

import (
	"context"

	"github.com/umk/jsonrpc2"
)

func GetCompletionRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetCompletionRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Completion(ctx, req.Messages, req.Params)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.Response(GetCompletionResponse{
		Completion: resp,
	})
}

func GetEmbeddingsRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetEmbeddingsRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Embeddings(ctx, req.Input, req.Params)
	if err != nil {
		return nil, newEmbeddingsError(err)
	}

	return c.Response(GetEmbeddingsResponse{
		Embeddings: resp,
	})
}

func GetStatisticsRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetStatisticsRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	resp := GetStatisticsResponse{
		BytesPerTok: cl.Samples.BytesPerTok(),
	}

	return c.Response(resp)
}
