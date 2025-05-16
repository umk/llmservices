package openai

import (
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/internal/pointer"
	"github.com/umk/llmservices/pkg/adapter"
)

func TestGetEmbeddingsRequest(t *testing.T) {
	input := "test input"
	params := adapter.EmbeddingsParams{
		Dimensions: pointer.Ptr(int64(128)),
		Model:      "text-embedding-ada-002",
	}

	expected := openai.EmbeddingNewParams{
		Dimensions: openai.Int(128),
		Model:      "text-embedding-ada-002",
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String("test input"),
		},
	}

	result := getEmbeddingsParams(input, params)
	assert.Equal(t, expected, result)
}

func TestGetEmbeddingsResponse(t *testing.T) {
	resp := &openai.CreateEmbeddingResponse{
		Data: []openai.Embedding{
			{
				Embedding: []float64{0.1, 0.2, 0.3},
			},
		},
		Usage: openai.CreateEmbeddingResponseUsage{
			PromptTokens: 10,
		},
	}

	expected := adapter.Embeddings{
		Data: []float64{0.1, 0.2, 0.3},
		Usage: &adapter.EmbeddingsUsage{
			PromptTokens: 10,
		},
	}

	result := getEmbeddingsResponse(resp)
	assert.Equal(t, expected, result)
}
