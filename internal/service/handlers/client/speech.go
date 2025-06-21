package client

import (
	"context"

	"github.com/umk/jsonrpc2"
)

func GetSpeechRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getSpeechRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(req.ClientID)
	if err != nil {
		return nil, err
	}

	resp, err := cl.Speech(ctx, req.Message, req.Params)
	if err != nil {
		return nil, newSpeechError(err)
	}

	return c.GetResponse(getSpeechResponse{
		Speech: resp,
	})
}
