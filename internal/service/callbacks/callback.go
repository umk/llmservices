package callbacks

import (
	"context"

	"github.com/umk/llmservices/pkg/adapter"
)

// Callback is a unified callback interface that aggregates callbacks of
// individual packages like threads and agents.
type Callback struct{}

func (c Callback) Call(ctx context.Context, fn adapter.ToolCallFunction) (string, error) {
	var res GetFunctionCallResponse
	if err := GetFunctionCallRPC(ctx, GetFunctionCallRequest{
		ToolCallFunction: fn,
	}, &res); err != nil {
		return "", err
	}

	return res.Response, nil
}

func (c Callback) Thought(ctx context.Context, content string) error {
	return PushThoughtRPC(ctx, PushThoughtRequest{
		Content: content,
	})
}
