package thread

import (
	"context"
	"errors"

	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client"
)

type Completion struct {
	// The same thread passed to the request.
	Thread Thread                   `json:"thread" validate:"required"`
	Usage  *adapter.CompletionUsage `json:"usage,omitempty"`
}

func (c *Client) Completion(ctx context.Context, thread Thread, params adapter.CompletionParams) (Completion, error) {
	if len(thread.Frames) == 0 {
		return Completion{}, errors.New("thread must have at least one frame")
	}

	var m []adapter.Message
	for _, f := range thread.Frames {
		m = append(m, f.Messages...)
	}

	resp, err := (*client.Client)(c).Completion(ctx, m, params)
	if err != nil {
		return Completion{}, err
	}

	f := &thread.Frames[len(thread.Frames)-1]

	message := resp.Message
	f.Messages = append(f.Messages, adapter.Message{
		OfAssistantMessage: &message,
	})
	f.Tokens = resp.Usage.PromptTokens + resp.Usage.CompletionTokens

	// Assign the token counts to frames after client stats have been updated.
	SetFrameTokens(&thread, c.Samples)

	return Completion{
		Thread: thread,
		Usage:  resp.Usage,
	}, nil
}
