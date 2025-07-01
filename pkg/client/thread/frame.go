package thread

import (
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client"
)

type MessagesFrame struct {
	Messages []adapter.Message `json:"messages" validate:"required,dive"`

	// FrameTokens in the frame. Derived from total tokens
	FrameTokens int64 `json:"tokens"`

	// The total tokens of the current and previous frames.
	Tokens int64 `json:"total_tokens"`
}

func (f *MessagesFrame) First() (*adapter.Message, bool) {
	if n := len(f.Messages); n > 0 {
		return &f.Messages[0], true
	}

	return nil, false
}

func (f *MessagesFrame) Last() (*adapter.Message, bool) {
	if n := len(f.Messages); n > 0 {
		return &f.Messages[n-1], true
	}

	return nil, false
}

func getEstimatedFrameSize(frame *MessagesFrame) int64 {
	var size int64
	for _, m := range frame.Messages {
		size += getEstimatedMessageSize(&m)
	}
	return size
}

func getEstimatedMessageSize(message *adapter.Message) int64 {
	var size int64

	switch {
	case message.OfSystemMessage != nil:
		size += int64(len(message.OfSystemMessage.Content))

	case message.OfUserMessage != nil:
		for _, p := range message.OfUserMessage.Parts {
			switch {
			case p.OfContentPartText != nil:
				size += int64(len(p.OfContentPartText.Text))
			case p.OfContentPartImageUrl != nil:
				size += int64(len(p.OfContentPartImageUrl.ImageUrl))
			}
		}

	case message.OfToolMessage != nil:
		size += int64(len(message.OfToolMessage.ToolCallID))
		for _, p := range message.OfToolMessage.Content {
			size += int64(len(p.Text))
		}

	case message.OfAssistantMessage != nil:
		switch {
		case message.OfAssistantMessage.Content != nil:
			size += int64(len(*message.OfAssistantMessage.Content))
		case message.OfAssistantMessage.Refusal != nil:
			size += int64(len(*message.OfAssistantMessage.Refusal))
		}

		for _, c := range message.OfAssistantMessage.ToolCalls {
			size += int64(len(c.ID))
			size += int64(len(c.Function.Name))
			size += int64(len(c.Function.Arguments))
		}
	}

	return size
}

// SetFrameTokens calculates and assigns the number of tokens for each frame in the thread.
func SetFrameTokens(thread *Thread, samples *client.Samples) {
	b := samples.BytesPerTok()

	var tokens int64
	for i := range thread.Frames {
		f := &thread.Frames[i]
		if len(f.Messages) == 0 {
			continue
		}
		if f.Tokens > 0 {
			d := f.Tokens - tokens
			f.FrameTokens = max(d, 0)
			tokens = f.Tokens
		} else {
			size := getEstimatedFrameSize(f)
			f.FrameTokens = int64(float32(size) / b)
			tokens += f.FrameTokens
		}
	}
}
