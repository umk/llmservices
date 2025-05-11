package openai

import (
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/internal/pointer"
	"github.com/umk/llmservices/pkg/adapter"
)

func TestGetCompletionRequest(t *testing.T) {
	// Test with minimal required fields
	t.Run("Basic Request", func(t *testing.T) {
		req := &adapter.CompletionRequest{
			Model: "gpt-4",
			Messages: []adapter.Message{
				adapter.CreateSystemMessage("You are a helpful assistant"),
				adapter.CreateUserMessage("Hello, world!"),
			},
		}

		params := getCompletionRequest(req)

		assert.Equal(t, "gpt-4", params.Model)
		assert.Len(t, params.Messages, 2)
		assert.NotNil(t, params.Messages[0].OfSystem)
		assert.NotNil(t, params.Messages[1].OfUser)
	})

	// Test with all fields populated
	t.Run("Full Request", func(t *testing.T) {
		freqPenalty := 0.5
		presPenalty := 0.7
		temp := 0.8
		topP := 0.9
		strict := true

		responseFormat := &adapter.ResponseFormat{
			OfResponseFormatText: &adapter.ResponseFormatText{},
		}

		req := &adapter.CompletionRequest{
			Model: "gpt-4",
			Messages: []adapter.Message{
				adapter.CreateSystemMessage("You are a helpful assistant"),
			},
			FrequencyPenalty: &freqPenalty,
			PresencePenalty:  &presPenalty,
			Temperature:      &temp,
			TopP:             &topP,
			ResponseFormat:   responseFormat,
			Stop:             []string{"STOP", "END"},
			Tools: []adapter.Tool{
				{
					Function: adapter.ToolFunction{
						Name:        "getWeather",
						Description: pointer.Ptr("Get the weather"),
						Parameters: map[string]any{
							"type": "object",
							"properties": map[string]any{
								"location": map[string]any{
									"type": "string",
								},
							},
						},
						Strict: &strict,
					},
				},
			},
		}

		params := getCompletionRequest(req)

		assert.Equal(t, "gpt-4", params.Model)
		// Compare Opt values properly
		assert.Equal(t, openai.Float(freqPenalty), params.FrequencyPenalty)
		assert.Equal(t, openai.Float(presPenalty), params.PresencePenalty)
		assert.Equal(t, openai.Float(temp), params.Temperature)
		assert.Equal(t, openai.Float(topP), params.TopP)
		assert.NotNil(t, params.ResponseFormat.OfText)
		assert.Equal(t, []string{"STOP", "END"}, params.Stop.OfChatCompletionNewsStopArray)
		assert.Len(t, params.Tools, 1)
		assert.Equal(t, "getWeather", params.Tools[0].Function.Name)
		assert.Equal(t, openai.String("Get the weather"), params.Tools[0].Function.Description)
		assert.Equal(t, openai.Bool(strict), params.Tools[0].Function.Strict)
	})

	t.Run("JSON Schema Response Format", func(t *testing.T) {
		description := "A weather response"
		strict := false
		responseFormat := &adapter.ResponseFormat{
			OfResponseFormatJSONSchema: &adapter.ResponseFormatJSONSchema{
				JSONSchema: adapter.JSONSchema{
					Name:        "WeatherResponse",
					Description: &description,
					Schema: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"temperature": map[string]any{
								"type": "number",
							},
						},
					},
					Strict: &strict,
				},
			},
		}

		req := &adapter.CompletionRequest{
			Model:          "gpt-4",
			Messages:       []adapter.Message{adapter.CreateUserMessage("What's the weather?")},
			ResponseFormat: responseFormat,
		}

		params := getCompletionRequest(req)

		assert.NotNil(t, params.ResponseFormat.OfJSONSchema)
		assert.Equal(t, "WeatherResponse", params.ResponseFormat.OfJSONSchema.JSONSchema.Name)
		assert.Equal(t, openai.String(description), params.ResponseFormat.OfJSONSchema.JSONSchema.Description)
		assert.Equal(t, openai.Bool(strict), params.ResponseFormat.OfJSONSchema.JSONSchema.Strict)
	})
}

func TestGetCompletionResponse(t *testing.T) {
	t.Run("Content Response", func(t *testing.T) {
		content := "This is the assistant's response"
		resp := &openai.ChatCompletion{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: content,
					},
				},
			},
			Usage: openai.CompletionUsage{
				PromptTokens:     10,
				CompletionTokens: 5,
			},
		}

		result, err := getCompletionResponse(resp)

		assert.NoError(t, err)
		assert.NotNil(t, result.Message.Content)
		assert.Equal(t, content, *result.Message.Content)
		assert.Nil(t, result.Message.Refusal)
		assert.Empty(t, result.Message.ToolCalls)
		assert.Equal(t, int64(10), result.Usage.PromptTokens)
		assert.Equal(t, int64(5), result.Usage.CompletionTokens)
	})

	t.Run("Refusal Response", func(t *testing.T) {
		refusal := "I cannot comply with that request"
		resp := &openai.ChatCompletion{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Refusal: refusal,
					},
				},
			},
			Usage: openai.CompletionUsage{
				PromptTokens:     10,
				CompletionTokens: 7,
			},
		}

		result, err := getCompletionResponse(resp)

		assert.NoError(t, err)
		assert.NotNil(t, result.Message.Refusal)
		assert.Equal(t, refusal, *result.Message.Refusal)
		assert.Nil(t, result.Message.Content)
		assert.Empty(t, result.Message.ToolCalls)
	})

	t.Run("Tool Call Response", func(t *testing.T) {
		resp := &openai.ChatCompletion{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						ToolCalls: []openai.ChatCompletionMessageToolCall{
							{
								ID: "call_123",
								Function: openai.ChatCompletionMessageToolCallFunction{
									Name:      "getWeather",
									Arguments: `{"location":"New York"}`,
								},
							},
						},
					},
				},
			},
			Usage: openai.CompletionUsage{
				PromptTokens:     15,
				CompletionTokens: 10,
			},
		}

		result, err := getCompletionResponse(resp)

		assert.NoError(t, err)
		assert.Len(t, result.Message.ToolCalls, 1)
		assert.Equal(t, "call_123", result.Message.ToolCalls[0].Id)
		assert.Equal(t, "getWeather", result.Message.ToolCalls[0].Function.Name)
		assert.Equal(t, `{"location":"New York"}`, result.Message.ToolCalls[0].Function.Arguments)
	})

	t.Run("Error - No Choices", func(t *testing.T) {
		resp := &openai.ChatCompletion{
			Choices: []openai.ChatCompletionChoice{},
		}

		_, err := getCompletionResponse(resp)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected number of choices")
	})

	t.Run("Error - Multiple Choices", func(t *testing.T) {
		resp := &openai.ChatCompletion{
			Choices: []openai.ChatCompletionChoice{
				{Message: openai.ChatCompletionMessage{Content: "Response 1"}},
				{Message: openai.ChatCompletionMessage{Content: "Response 2"}},
			},
		}

		_, err := getCompletionResponse(resp)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected number of choices")
	})
}
