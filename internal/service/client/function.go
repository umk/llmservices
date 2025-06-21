package client

import (
	"context"

	"github.com/umk/jsonrpc2"
	"github.com/umk/llmservices/pkg/adapter"
)

type FunctionCaller struct {
	Client jsonrpc2.Client
}

func (c *FunctionCaller) Call(ctx context.Context, fn adapter.ToolCallFunction) (string, error) {
	var resp functionCallResponse
	if err := c.Client.Call(ctx, "getFunctionCall", functionCallRequest{
		ToolCallFunction: fn,
	}, &resp); err != nil {
		return "", err
	}

	return resp.Response, nil
}
