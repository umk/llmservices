package adapter

type Message struct {
	OfSystemMessage    *SystemMessage    `json:"system,omitempty"`
	OfUserMessage      *UserMessage      `json:"user,omitempty"`
	OfAssistantMessage *AssistantMessage `json:"assistant,omitempty"`
	OfToolMessage      *ToolMessage      `json:"tool,omitempty"`
}

type SystemMessage struct {
	Content string `json:"content" validate:"required"`
}

type UserMessage struct {
	Parts []ContentPart `json:"parts" validate:"required,min=1"`
}

type AssistantMessage struct {
	Content   *string    `json:"content,omitempty"`
	Refusal   *string    `json:"refusal,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolMessage struct {
	ToolCallId string            `json:"tool_call_id" validate:"required"`
	Content    []ContentPartText `json:"content" validate:"required,min=1"`
}

func CreateUserMessage(content string) Message {
	return Message{
		OfUserMessage: &UserMessage{
			Parts: []ContentPart{{
				OfContentPartText: &ContentPartText{
					Text: content,
				},
			}},
		},
	}
}

func CreateSystemMessage(content string) Message {
	return Message{
		OfSystemMessage: &SystemMessage{
			Content: content,
		},
	}
}

func CreateToolMessage(callId string, response string) Message {
	return Message{
		OfToolMessage: &ToolMessage{
			ToolCallId: callId,
			Content: []ContentPartText{{
				Text: response,
			}},
		},
	}
}
