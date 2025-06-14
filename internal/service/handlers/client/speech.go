package client

import (
	"context"

	"github.com/umk/jsonrpc2"
)

func GetSpeech(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req getSpeechRequest
	if err := c.GetRequestBody(&req); err != nil {
		return nil, err
	}

	cl := getClient(req.ClientID)
	if cl == nil {
		return nil, errClientNotFound
	}

	resp, err := cl.Speech(ctx, req.Message, req.Params)
	if err != nil {
		return nil, newSpeechError(err)
	}

	return c.GetResponse(getSpeechResponse{
		Speech: resp,
	})
}
