package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/pkg/adapter"
)

func TestGetMessage_SystemMessage(t *testing.T) {
	message := &adapter.Message{
		OfSystemMessage: &adapter.SystemMessage{
			Content: "This is a system message",
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfSystem)
	// Can't directly access Opt value, so we just check it exists
	assert.NotNil(t, result.OfSystem.Content.OfString)
}

func TestGetMessage_UserMessage(t *testing.T) {
	message := &adapter.Message{
		OfUserMessage: &adapter.UserMessage{
			Parts: []adapter.ContentPart{
				{
					OfContentPartText: &adapter.ContentPartText{
						Text: "This is a user message",
					},
				},
			},
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfUser)
	assert.Len(t, result.OfUser.Content.OfArrayOfContentParts, 1)
	assert.Equal(t, "This is a user message", result.OfUser.Content.OfArrayOfContentParts[0].OfText.Text)
}

func TestGetMessage_UserMessageMultipleParts(t *testing.T) {
	message := &adapter.Message{
		OfUserMessage: &adapter.UserMessage{
			Parts: []adapter.ContentPart{
				{
					OfContentPartText: &adapter.ContentPartText{
						Text: "Part 1",
					},
				},
				{
					OfContentPartImageUrl: &adapter.ContentPartImage{
						ImageUrl: "https://example.com/image.jpg",
					},
				},
			},
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfUser)
	assert.Len(t, result.OfUser.Content.OfArrayOfContentParts, 2)
	assert.Equal(t, "Part 1", result.OfUser.Content.OfArrayOfContentParts[0].OfText.Text)
	assert.Equal(t, "https://example.com/image.jpg", result.OfUser.Content.OfArrayOfContentParts[1].OfImageURL.ImageURL.URL)
}

func TestGetMessage_AssistantMessageWithContent(t *testing.T) {
	content := "This is an assistant response"
	message := &adapter.Message{
		OfAssistantMessage: &adapter.AssistantMessage{
			Content: &content,
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfAssistant)
	// We can't directly access the Opt value, so we just check it exists
	assert.NotNil(t, result.OfAssistant.Content.OfString)
}

func TestGetMessage_AssistantMessageWithRefusal(t *testing.T) {
	refusal := "I cannot comply with that request"
	message := &adapter.Message{
		OfAssistantMessage: &adapter.AssistantMessage{
			Refusal: &refusal,
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfAssistant)
	// We can't directly access the Opt value, so we just check it exists
	assert.NotNil(t, result.OfAssistant.Refusal)
}

func TestGetMessage_AssistantMessageWithToolCalls(t *testing.T) {
	message := &adapter.Message{
		OfAssistantMessage: &adapter.AssistantMessage{
			ToolCalls: []adapter.ToolCall{
				{
					ID: "call_123",
					Function: adapter.ToolCallFunction{
						Name:      "get_weather",
						Arguments: `{"location":"New York"}`,
					},
				},
			},
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfAssistant)
	assert.Len(t, result.OfAssistant.ToolCalls, 1)
	assert.Equal(t, "call_123", result.OfAssistant.ToolCalls[0].ID)
	assert.Equal(t, "get_weather", result.OfAssistant.ToolCalls[0].Function.Name)
	assert.Equal(t, `{"location":"New York"}`, result.OfAssistant.ToolCalls[0].Function.Arguments)
}

func TestGetMessage_ToolMessage(t *testing.T) {
	message := &adapter.Message{
		OfToolMessage: &adapter.ToolMessage{
			ToolCallID: "call_123",
			Content: []adapter.ContentPartText{
				{
					Text: "Tool response",
				},
			},
		},
	}

	result := getMessage(message)

	assert.NotNil(t, result.OfTool)
	assert.Equal(t, "call_123", result.OfTool.ToolCallID)
	assert.Len(t, result.OfTool.Content.OfArrayOfContentParts, 1)
	assert.Equal(t, "Tool response", result.OfTool.Content.OfArrayOfContentParts[0].Text)
}

func TestGetMessage_EmptyMessage(t *testing.T) {
	message := &adapter.Message{}

	result := getMessage(message)

	// Check fields are nil
	assert.Nil(t, result.OfSystem)
	assert.Nil(t, result.OfUser)
	assert.Nil(t, result.OfAssistant)
	assert.Nil(t, result.OfTool)
}
