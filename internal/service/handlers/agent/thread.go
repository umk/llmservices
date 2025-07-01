package agent

import (
	"context"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/internal/service/callbacks"
)

func GetResponseRPC(ctx context.Context, c jsonrpc2.RPCContext) (any, error) {
	var req GetResponseRequest
	if err := c.Request(&req); err != nil {
		return nil, err
	}

	cl, err := GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}

	req.Params.Handler = callbacks.Callback{}

	resp, err := cl.Response(ctx, req.Thread, req.Params)
	if err != nil {
		return nil, newResponseError(err)
	}

	return c.Response(GetResponseResponse{
		Response: resp,
	})
}
