package adapter

import "context"

type Adapter interface {
	GetCompletion(ctx context.Context, req *CompletionRequest) (CompletionResponse, error)
	GetEmbeddings(ctx context.Context, req *EmbeddingsRequest) (EmbeddingsResponse, error)
}
