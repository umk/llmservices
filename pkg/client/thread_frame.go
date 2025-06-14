package client

import (
	"errors"

	"github.com/umk/llmservices/pkg/adapter"
)

type MessagesFrame struct {
	Messages []adapter.Message `json:"messages" validate:"required,dive"`

	// FrameTokens in the frame. Derived from total tokens
	FrameTokens int64 `json:"tokens"`

	// The total tokens of the current and previous frames.
	Tokens int64 `json:"total_tokens"`
}

func (f *MessagesFrame) Response() (*adapter.AssistantMessage, error) {
	if n := len(f.Messages); n > 0 {
		m := f.Messages[n-1]
		if m.OfAssistantMessage != nil {
			return m.OfAssistantMessage, nil
		}
	}

	return nil, errors.New("frame doesn't contain an assistant message")
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
