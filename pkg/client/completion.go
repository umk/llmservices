package client

import (
	"context"

	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Client) GetCompletion(ctx context.Context, req *adapter.CompletionRequest) (
	adapter.CompletionResponse, error,
) {
	if err := c.s.Acquire(ctx, 1); err != nil {
		return adapter.CompletionResponse{}, err
	}
	defer c.s.Release(1)

	resp, err := c.adapter.GetCompletion(ctx, req)

	if err == nil {
		c.setSamplesFromCompl(&resp)
	}

	return resp, err
}

func (c *Client) setSamplesFromCompl(resp *adapter.CompletionResponse) {
	toks := resp.Usage.CompletionTokens
	if toks == 0 {
		return
	}

	if len(resp.Message.ToolCalls) > 0 || resp.Message.Refusal != nil || resp.Message.Content == nil {
		return
	}

	if b := len(*resp.Message.Content); b >= minSampleSize {
		c.Samples.put(float32(b) / float32(toks))
	}
}
