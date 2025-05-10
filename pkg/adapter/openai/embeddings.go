package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Adapter) GetEmbeddings(ctx context.Context, req *adapter.EmbeddingsRequest) (adapter.EmbeddingsResponse, error) {
	params := getEmbeddingsRequest(req)

	resp, err := c.Client.Embeddings.New(ctx, params)
	if err != nil {
		return adapter.EmbeddingsResponse{}, err
	}
	if len(resp.Data) != 1 {
		return adapter.EmbeddingsResponse{}, fmt.Errorf("unexpected number of embeddings: %d", len(resp.Data))
	}

	return getEmbeddingsResponse(resp), nil
}

func getEmbeddingsRequest(req *adapter.EmbeddingsRequest) openai.EmbeddingNewParams {
	return openai.EmbeddingNewParams{
		Dimensions: getOpt(req.Dimensions),
		Model:      req.Model,
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(req.Input),
		},
	}
}

func getEmbeddingsResponse(resp *openai.CreateEmbeddingResponse) adapter.EmbeddingsResponse {
	return adapter.EmbeddingsResponse{
		Data: resp.Data[0].Embedding,
		Usage: &adapter.EmbeddingsUsage{
			PromptTokens: resp.Usage.PromptTokens,
		},
	}
}
