package thread

import (
	"context"
	"errors"

	"github.com/umk/llmservices/internal/msg"
	"github.com/umk/llmservices/pkg/adapter"
)

type Summarizer struct {
	client *Client

	fraction float32 // fraction of history to be summarized

	maxTokens   *int
	maxMessages *int
}

type SummarizerOption func(*Summarizer)

func WithMaxTokens(maxTokens int) SummarizerOption {
	return func(s *Summarizer) {
		s.maxTokens = &maxTokens
	}
}

func WithMaxMessages(maxMessages int) SummarizerOption {
	return func(s *Summarizer) {
		s.maxMessages = &maxMessages
	}
}

func NewSummarizer(client *Client, fraction float32, opts ...SummarizerOption) *Summarizer {
	const (
		fractionMin = 0.1
		fractionMax = 1
	)

	if fraction > fractionMax {
		fraction = fractionMax
	} else if fraction < fractionMin {
		fraction = fractionMin
	}

	s := &Summarizer{
		client:   client,
		fraction: fraction,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Summarizer) Summarize(ctx context.Context, thread Thread) (Thread, error) {
	if !s.checkConditions(&thread) {
		return thread, nil
	}

	n := int(float32(len(thread.Frames)) * s.fraction)

	if n < 2 {
		return thread, nil
	}

	v, err := s.getSummary(ctx, thread.Frames[:n])
	if err != nil {
		return Thread{}, err
	}

	m, err := msg.RenderSummaryMessage(msg.SummaryMessageParams{
		Summary: v,
	})
	if err != nil {
		return Thread{}, err
	}

	var t Thread
	if s, ok := getSystemFrame(thread); ok {
		t.Frames = append(t.Frames, s)
	}
	t.Frames = append(t.Frames, MessagesFrame{
		Messages: []adapter.Message{adapter.CreateAssistantMessage(m)},
	})
	t.Frames = append(t.Frames, thread.Frames[n:]...)

	return t, nil
}

func getSystemFrame(thread Thread) (MessagesFrame, bool) {
	if len(thread.Frames) == 0 {
		return MessagesFrame{}, false
	}

	f := thread.Frames[0]
	if len(f.Messages) == 0 {
		return MessagesFrame{}, false
	}

	m := f.Messages[0]
	if m.OfSystemMessage == nil {
		return MessagesFrame{}, false
	}

	if len(f.Messages) == 1 {
		return f, true
	}

	return MessagesFrame{
		Messages: []adapter.Message{m},
	}, true
}

func (s *Summarizer) checkConditions(thread *Thread) bool {
	if s.maxTokens != nil {
		t := thread.Tokens(s.client.Samples)
		if t >= int64(*s.maxTokens) {
			return true
		}
	}

	if s.maxMessages != nil {
		m := 0
		for _, f := range thread.Frames {
			m += len(f.Messages)
		}

		if m >= *s.maxMessages {
			return true
		}
	}

	return false
}

func (s *Summarizer) getSummary(ctx context.Context, frames []MessagesFrame) (string, error) {
	m, err := msg.RenderSummarizeMessage(msg.SummarizeMessageParams{})
	if err != nil {
		return "", err
	}

	t := Thread{
		Frames: make([]MessagesFrame, len(frames), len(frames)+1),
	}

	copy(t.Frames, frames)

	t.Frames = append(t.Frames, MessagesFrame{
		Messages: []adapter.Message{adapter.CreateUserMessage(m)},
	})

	c, err := s.client.Completion(ctx, t, adapter.CompletionParams{})
	if err != nil {
		return "", err
	}

	r, err := c.Thread.Response()
	if err != nil {
		return "", err
	}

	if r.Refusal != nil {
		return "", errors.New("generating summary was refused")
	}

	return r.Text()
}
