package client

import (
	"context"

	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Client) GetEmbeddings(ctx context.Context, req *adapter.EmbeddingsRequest) (
	adapter.EmbeddingsResponse, error,
) {
	if err := c.s.Acquire(ctx, 1); err != nil {
		return adapter.EmbeddingsResponse{}, err
	}
	defer c.s.Release(1)

	resp, err := c.adapter.GetEmbeddings(ctx, req)

	if err == nil {
		c.setSamplesFromEmbedding(req, &resp)
	}

	return resp, err
}

func (c *Client) setSamplesFromEmbedding(req *adapter.EmbeddingsRequest, resp *adapter.EmbeddingsResponse) {
	toks := resp.Usage.PromptTokens
	if toks == 0 {
		return
	}

	if b := len(req.Input); b >= minSampleSize {
		c.Samples.put(float32(b) / float32(toks))
	}
}
