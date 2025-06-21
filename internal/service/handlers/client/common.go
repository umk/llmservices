package client

import (
	"context"

	"github.com/umk/jsonrpc2"
)

func GetCompletionRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getCompletionRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Completion(ctx, req.Messages, req.Params)
	if err != nil {
		return nil, newCompletionError(err)
	}

	return c.GetResponse(getCompletionResponse{
		Completion: resp,
	})
}

func GetEmbeddingsRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getEmbeddingsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Embeddings(ctx, req.Input, req.Params)
	if err != nil {
		return nil, newEmbeddingsError(err)
	}

	return c.GetResponse(getEmbeddingsResponse{
		Embeddings: resp,
	})
}

func GetStatisticsRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getStatisticsRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	resp := getStatisticsResponse{
		BytesPerTok: cl.Samples.BytesPerTok(),
	}

	return c.GetResponse(resp)
}
