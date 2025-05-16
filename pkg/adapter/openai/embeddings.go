package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Adapter) Embeddings(ctx context.Context, input string, params adapter.EmbeddingsParams) (adapter.Embeddings, error) {
	p := getEmbeddingsParams(input, params)

	resp, err := c.Client.Embeddings.New(ctx, p)
	if err != nil {
		return adapter.Embeddings{}, err
	}
	if len(resp.Data) != 1 {
		return adapter.Embeddings{}, fmt.Errorf("unexpected number of embeddings: %d", len(resp.Data))
	}

	return getEmbeddingsResponse(resp), nil
}

func getEmbeddingsParams(input string, params adapter.EmbeddingsParams) openai.EmbeddingNewParams {
	return openai.EmbeddingNewParams{
		Dimensions: getOpt(params.Dimensions),
		Model:      params.Model,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(input),
		},
	}
}

func getEmbeddingsResponse(resp *openai.CreateEmbeddingResponse) adapter.Embeddings {
	return adapter.Embeddings{
		Data: resp.Data[0].Embedding,
		Usage: &adapter.EmbeddingsUsage{
			PromptTokens: resp.Usage.PromptTokens,
		},
	}
}
