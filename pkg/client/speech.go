package client

import (
	"context"

	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Client) Speech(ctx context.Context, message adapter.SpeechMessage, params adapter.SpeechParams) (
	adapter.Speech, error,
) {
	if err := c.s.Acquire(ctx, 1); err != nil {
		return adapter.Speech{}, err
	}
	defer c.s.Release(1)

	a, ok := c.adapter.(adapter.SpeechAdapter)
	if !ok {
		return adapter.Speech{}, ErrNotSupportedByAdapter
	}

	return a.Speech(ctx, message, params)
}
