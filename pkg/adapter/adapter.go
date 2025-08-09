package adapter

import "context"

type Adapter interface {
	Completion(ctx context.Context, messages []Message, params CompletionParams) (Completion, error)
	Embeddings(ctx context.Context, input string, params EmbeddingsParams) (Embeddings, error)
}
