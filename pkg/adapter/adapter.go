package adapter

import "context"

type Adapter interface {
	Completion(ctx context.Context, messages []Message, params CompletionParams) (Completion, error)
	Embeddings(ctx context.Context, input string, params EmbeddingsParams) (Embeddings, error)
}

type SpeechAdapter interface {
	Speech(ctx context.Context, message SpeechMessage, params SpeechParams) (Speech, error)
}
