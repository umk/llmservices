package handlers

import (
	"context"

	"github.com/umk/jsonrpc2"
)

func GetSpeechRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetSpeechRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Speech(ctx, req.Message, req.Params)
	if err != nil {
		return nil, newSpeechError(err)
	}

	return c.Response(GetSpeechResponse{
		Speech: resp,
	})
}
