package client

import (
	"context"

	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Client) Completion(ctx context.Context, messages []adapter.Message, params adapter.CompletionParams) (
	adapter.Completion, error,
) {
	if err := c.s.Acquire(ctx, 1); err != nil {
		return adapter.Completion{}, err
	}
	defer c.s.Release(1)

	// If the model is not set, use the default one
	if params.Model == "" {
		params.Model = c.config.Model
	}

	resp, err := c.adapter.Completion(ctx, messages, params)

	if err == nil {
		c.setSamplesFromCompl(&resp)
	}

	return resp, err
}

func (c *Client) setSamplesFromCompl(resp *adapter.Completion) {
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
