package thread

import (
	"errors"

	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client"
)

type Thread struct {
	Frames []MessagesFrame `json:"frames" validate:"required,dive,min=1"`
}

func (t *Thread) Request() (adapter.UserMessage, error) {
	if m, ok := t.Last(); ok && m.OfUserMessage != nil {
		return *m.OfUserMessage, nil
	}

	return adapter.UserMessage{}, errors.New("frame doesn't contain a user message")
}

func (t *Thread) Response() (adapter.AssistantMessage, error) {
	if m, ok := t.Last(); ok && m.OfAssistantMessage != nil {
		return *m.OfAssistantMessage, nil
	}

	return adapter.AssistantMessage{}, errors.New("frame doesn't contain an assistant message")
}

func (t *Thread) First() (*adapter.Message, bool) {
	for i := range len(t.Frames) {
		if m, ok := t.Frames[i].First(); ok {
			return m, true
		}
	}

	return nil, false
}

func (t *Thread) Last() (*adapter.Message, bool) {
	for i := len(t.Frames) - 1; i >= 0; i-- {
		if m, ok := t.Frames[i].Last(); ok {
			return m, true
		}
	}

	return nil, false
}

func (t *Thread) Tokens(samples *client.Samples) int64 {
	var toks int64

	i := len(t.Frames)
	for ; i > 0; i-- {
		f := t.Frames[i-1]
		if f.Tokens > 0 {
			toks = f.Tokens
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
