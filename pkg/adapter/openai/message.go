package openai

import (
	"github.com/openai/openai-go"
	"github.com/umk/llmservices/pkg/adapter"
)

func getMessage(message *adapter.Message) openai.ChatCompletionMessageParamUnion {
	switch {
	case message.OfSystemMessage != nil:
		return getSystemMessage(message.OfSystemMessage)
	case message.OfUserMessage != nil:
		return getUserMessage(message.OfUserMessage)
	case message.OfAssistantMessage != nil:
		return getAssistantMessage(message.OfAssistantMessage)
	case message.OfToolMessage != nil:
		return getToolMessage(message.OfToolMessage)
	default:
		return openai.ChatCompletionMessageParamUnion{}
	}
}

func getSystemMessage(systemMessage *adapter.SystemMessage) openai.ChatCompletionMessageParamUnion {
	return openai.ChatCompletionMessageParamUnion{
		OfSystem: &openai.ChatCompletionSystemMessageParam{
			Content: openai.ChatCompletionSystemMessageParamContentUnion{
				OfString: openai.String(systemMessage.Content),
			},
		},
	}
}

func getUserMessage(userMessage *adapter.UserMessage) openai.ChatCompletionMessageParamUnion {
	result := openai.ChatCompletionMessageParamUnion{
		OfUser: &openai.ChatCompletionUserMessageParam{},
	}
	for _, part := range userMessage.Parts {
		result.OfUser.Content.OfArrayOfContentParts = append(result.OfUser.Content.OfArrayOfContentParts, getContentPart(&part))
	}
	return result
}

func getAssistantMessage(assistantMessage *adapter.AssistantMessage) openai.ChatCompletionMessageParamUnion {
	result := openai.ChatCompletionMessageParamUnion{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			Content: openai.ChatCompletionAssistantMessageParamContentUnion{
				OfString: getOpt(assistantMessage.Content),
			},
			Refusal: getOpt(assistantMessage.Refusal),
		},
	}
	for _, call := range assistantMessage.ToolCalls {
		result.OfAssistant.ToolCalls = append(result.OfAssistant.ToolCalls, openai.ChatCompletionMessageToolCallParam{
			ID: call.ID,
			Function: openai.ChatCompletionMessageToolCallFunctionParam{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments,
			},
		})
	}
	return result
}

func getToolMessage(toolMessage *adapter.ToolMessage) openai.ChatCompletionMessageParamUnion {
	result := openai.ChatCompletionMessageParamUnion{
		OfTool: &openai.ChatCompletionToolMessageParam{
			ToolCallID: toolMessage.ToolCallID,
		},
	}
	for _, part := range toolMessage.Content {
		result.OfTool.Content.OfArrayOfContentParts = append(
			result.OfTool.Content.OfArrayOfContentParts,
			openai.ChatCompletionContentPartTextParam{
				Text: part.Text,
			},
		)
	}
	return result
}
