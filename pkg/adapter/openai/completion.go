package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Adapter) GetCompletion(ctx context.Context, req *adapter.CompletionRequest) (adapter.CompletionResponse, error) {
	params := getCompletionRequest(req)

	resp, err := c.Client.Chat.Completions.New(ctx, params)
	if err != nil {
		return adapter.CompletionResponse{}, err
	}

	return getCompletionResponse(resp)
}

func getCompletionRequest(req *adapter.CompletionRequest) openai.ChatCompletionNewParams {
	params := openai.ChatCompletionNewParams{
		Model:            req.Model,
		FrequencyPenalty: getOpt(req.FrequencyPenalty),
		PresencePenalty:  getOpt(req.PresencePenalty),
		ResponseFormat:   getResponseFormat(req.ResponseFormat),
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfChatCompletionNewsStopArray: req.Stop,
		},
		Temperature: getOpt(req.Temperature),
		TopP:        getOpt(req.TopP),
	}
	for _, message := range req.Messages {
		params.Messages = append(params.Messages, getMessage(&message))
	}
	for _, tool := range req.Tools {
		params.Tools = append(params.Tools, getTool(&tool))
	}

	return params
}

func getCompletionResponse(resp *openai.ChatCompletion) (adapter.CompletionResponse, error) {
	if len(resp.Choices) != 1 {
		return adapter.CompletionResponse{}, fmt.Errorf("unexpected number of choices: %d", len(resp.Choices))
	}

	message := resp.Choices[0].Message
	var content, refusal *string
	if message.Refusal != "" {
		refusal = &message.Refusal
	} else if message.Content != "" {
		content = &message.Content
	}

	result := adapter.CompletionResponse{
		Message: adapter.AssistantMessage{
			Content: content,
			Refusal: refusal,
		},
		Usage: &adapter.CompletionUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
		},
	}

	for _, call := range message.ToolCalls {
		result.Message.ToolCalls = append(result.Message.ToolCalls, adapter.ToolCall{
			Id: call.ID,
			Function: adapter.ToolCallFunction{
				Name:      call.Function.Name,
				Arguments: call.Function.Arguments,
			},
		})
	}

	return result, nil
}

func getResponseFormat(format *adapter.ResponseFormat) openai.ChatCompletionNewParamsResponseFormatUnion {
	if format == nil {
		return openai.ChatCompletionNewParamsResponseFormatUnion{}
	}

	if format.OfResponseFormatText != nil {
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfText: &openai.ResponseFormatTextParam{},
		}
	}

	if format.OfResponseFormatJSONSchema != nil {
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        format.OfResponseFormatJSONSchema.JSONSchema.Name,
					Description: getOpt(format.OfResponseFormatJSONSchema.JSONSchema.Description),
					Schema:      format.OfResponseFormatJSONSchema.JSONSchema.Schema,
					Strict:      getOpt(format.OfResponseFormatJSONSchema.JSONSchema.Strict),
				},
			},
		}
	}

	return openai.ChatCompletionNewParamsResponseFormatUnion{}
}

func getTool(tool *adapter.Tool) openai.ChatCompletionToolParam {
	return openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        tool.Function.Name,
			Description: getOpt(tool.Function.Description),
			Parameters:  tool.Function.Parameters,
			Strict:      getOpt(tool.Function.Strict),
		},
	}
}
