package client

import (
	"context"
	"errors"

	"github.com/umk/llmservices/pkg/adapter"
)

type Thread struct {
	Frames []MessagesFrame `json:"frames" validate:"required,dive,min=1"`
}

type ThreadCompletion struct {
	// The same thread passed to the request.
	Thread *Thread                  `json:"thread" validate:"required"`
	Usage  *adapter.CompletionUsage `json:"usage,omitempty"`
}

func (t *Thread) Response() (*adapter.AssistantMessage, error) {
	if n := len(t.Frames); n > 0 {
		f := t.Frames[len(t.Frames)-1]
		return f.Response()
	}
	return nil, errors.New("thread must have at least one frame")
}

func (t *Thread) Tokens(samples *Samples) int64 {
	var toks int64

	i := len(t.Frames)
	for ; i > 0; i-- {
		f := t.Frames[i-1]
		if f.CompletionTokens > 0 {
			toks = f.CompletionTokens
			break
		}
	}

	var size int64
	for ; i < len(t.Frames); i++ {
		size += getEstimatedFrameSize(&t.Frames[i])
	}

	b := samples.BytesPerTok()
	return toks + int64(float32(size)/b)
}

func (c *Client) ThreadCompletion(ctx context.Context, thread *Thread, params adapter.CompletionParams) (ThreadCompletion, error) {
	if len(thread.Frames) == 0 {
		return ThreadCompletion{}, errors.New("thread must have at least one frame")
	}

	if err := c.s.Acquire(ctx, 1); err != nil {
		return ThreadCompletion{}, err
	}
	defer c.s.Release(1)

	// If the model is not set, use the default one
	if params.Model == "" {
		params.Model = c.config.Model
	}

	var m []adapter.Message
	for _, f := range thread.Frames {
		m = append(m, f.Messages...)
	}

	resp, err := c.adapter.Completion(ctx, m, params)
	if err != nil {
		return ThreadCompletion{}, err
	}

	c.setSamplesFromCompl(&resp)

	f := &thread.Frames[len(thread.Frames)-1]

	message := resp.Message
	f.Messages = append(f.Messages, adapter.Message{
		OfAssistantMessage: &message,
	})
	f.CompletionTokens = resp.Usage.CompletionTokens

	// Assign the token counts to frames after client stats have been updated.
	setFrameTokens(thread, c.Samples)

	return ThreadCompletion{
		Thread: thread,
		Usage:  resp.Usage,
	}, nil
}

func setFrameTokens(thread *Thread, samples *Samples) {
	b := samples.BytesPerTok()

	var completionTokens int64
	for i := range thread.Frames {
		f := &thread.Frames[i]
		if len(f.Messages) == 0 {
			continue
		}
		if f.CompletionTokens > 0 {
			d := f.CompletionTokens - completionTokens
			f.Tokens = max(d, 0)
			completionTokens = f.CompletionTokens
		} else {
			size := getEstimatedFrameSize(f)
			f.Tokens = int64(float32(size) / b)
			completionTokens += f.Tokens
		}
	}
}
