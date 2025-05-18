package client

import (
	"context"

	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Client) Embeddings(ctx context.Context, input string, params adapter.EmbeddingsParams) (
	adapter.Embeddings, error,
) {
	if err := c.s.Acquire(ctx, 1); err != nil {
		return adapter.Embeddings{}, err
	}
	defer c.s.Release(1)

	// If the model is not set, use the default one
	if params.Model == "" {
		params.Model = c.config.Model
	}

	resp, err := c.adapter.Embeddings(ctx, input, params)

	if err == nil {
		c.setSamplesFromEmbedding(input, &resp)
	}

	return resp, err
}

func (c *Client) setSamplesFromEmbedding(input string, resp *adapter.Embeddings) {
	toks := resp.Usage.PromptTokens
	if toks == 0 {
		return
	}

	if b := len(input); b >= minSampleSize {
		c.Samples.put(float32(b) / float32(toks))
	}
}
