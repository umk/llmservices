package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Adapter) Completion(ctx context.Context, messages []adapter.Message, params adapter.CompletionParams) (adapter.Completion, error) {
	p := getCompletionParams(messages, params)

	resp, err := c.Client.Chat.Completions.New(ctx, p)
	if err != nil {
		return adapter.Completion{}, err
	}

	return getCompletionResponse(resp)
}

func getCompletionParams(messages []adapter.Message, params adapter.CompletionParams) openai.ChatCompletionNewParams {
	r := openai.ChatCompletionNewParams{
		Model:            params.Model,
		FrequencyPenalty: getOpt(params.FrequencyPenalty),
		PresencePenalty:  getOpt(params.PresencePenalty),
		ResponseFormat:   getResponseFormat(params.ResponseFormat),
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfStringArray: params.Stop,
		},
		Temperature: getOpt(params.Temperature),
		TopP:        getOpt(params.TopP),
	}
	for _, message := range messages {
		r.Messages = append(r.Messages, getMessage(&message))
	}
	for _, tool := range params.Tools {
		r.Tools = append(r.Tools, getTool(&tool))
	}

	return r
}

func getCompletionResponse(resp *openai.ChatCompletion) (adapter.Completion, error) {
	if len(resp.Choices) != 1 {
		return adapter.Completion{}, fmt.Errorf("unexpected number of choices: %d", len(resp.Choices))
	}

	message := resp.Choices[0].Message
	var content, refusal *string
	if message.Refusal != "" {
		refusal = &message.Refusal
	} else if message.Content != "" {
		content = &message.Content
	}

	result := adapter.Completion{
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
			ID: call.ID,
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
