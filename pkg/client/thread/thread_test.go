package thread

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/pkg/adapter"
	"github.com/umk/llmservices/pkg/client"
)

func TestThread_Tokens(t *testing.T) {
	// Create a test samples object with a known BytesPerTok value
	samples := client.NewSamples(5, 4.0) // 4.0 bytes per token

	// Test cases
	tests := []struct {
		name           string
		thread         *Thread
		expectedTokens int64
	}{
		{
			name: "empty thread",
			thread: &Thread{
				Frames: []MessagesFrame{
					{
						Messages:    []adapter.Message{},
						FrameTokens: 0,
						Tokens:      0,
					},
				},
			},
			expectedTokens: 0,
		},
		{
			name: "thread with tokens already set",
			thread: &Thread{
				Frames: []MessagesFrame{
					{
						Messages: []adapter.Message{
							{
								OfSystemMessage: &adapter.SystemMessage{
									Content: "System message",
								},
							},
						},
						FrameTokens: 10,
						Tokens:      10,
					},
				},
			},
			expectedTokens: 10, // Should use the existing token count
		},
		{
			name: "thread with multiple frames, last has token count",
			thread: &Thread{
				Frames: []MessagesFrame{
					{
						Messages: []adapter.Message{
							{
								OfSystemMessage: &adapter.SystemMessage{
									Content: "System message",
								},
							},
						},
						FrameTokens: 10,
						Tokens:      10,
					},
					{
						Messages: []adapter.Message{
							{
								OfUserMessage: &adapter.UserMessage{
									Parts: []adapter.ContentPart{
										{
											OfContentPartText: &adapter.ContentPartText{
												Text: "Hello",
											},
										},
									},
								},
							},
						},
						FrameTokens: 0,
						Tokens:      20, // This is the latest token count
					},
				},
			},
			expectedTokens: 20, // Should use the last frame's token count
		},
		{
			name: "thread with estimation needed",
			thread: &Thread{
				Frames: []MessagesFrame{
					{
						Messages: []adapter.Message{
							{
								OfSystemMessage: &adapter.SystemMessage{
									Content: "System message",
								},
							},
						},
						FrameTokens: 10,
						Tokens:      10,
					},
					{
						Messages: []adapter.Message{
							{
								OfUserMessage: &adapter.UserMessage{
									Parts: []adapter.ContentPart{
										{
											OfContentPartText: &adapter.ContentPartText{
												Text: "Hello, this is a message that needs estimation", // 40 bytes
											},
										},
									},
								},
							},
						},
						// No tokens set, will need estimation
					},
				},
			},
			expectedTokens: 21, // 10 (from first frame) + 11 (40 bytes / 4 bytes per token, rounded up)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tt.thread.Tokens(samples)
			assert.Equal(t, tt.expectedTokens, tokens, "Token count mismatch")
		})
	}
}

func TestThread_TokensWithMultipleMessageTypes(t *testing.T) {
	// Create a test samples object with a known BytesPerTok value
	samples := client.NewSamples(5, 2.0) // 2.0 bytes per token

	// Create a content pointer for assistant message
	content := "I'm an assistant"

	// Create a thread with various message types
	thread := &Thread{
		Frames: []MessagesFrame{
			{
				Messages: []adapter.Message{
					{
						OfSystemMessage: &adapter.SystemMessage{
							Content: "System instructions", // 19 bytes
						},
					},
				},
				Tokens: 10, // Known token count
			},
			{
				Messages: []adapter.Message{
					{
						OfUserMessage: &adapter.UserMessage{
							Parts: []adapter.ContentPart{
								{
									OfContentPartText: &adapter.ContentPartText{
										Text: "User question", // 13 bytes
									},
								},
								{
									OfContentPartImageUrl: &adapter.ContentPartImage{
										ImageUrl: "https://example.com/image.jpg", // 29 bytes
									},
								},
							},
						},
					},
					{
						OfAssistantMessage: &adapter.AssistantMessage{
							Content: &content, // 16 bytes ("I'm an assistant")
							ToolCalls: []adapter.ToolCall{
								{
									ID: "tool1", // 5 bytes
									Function: adapter.ToolCallFunction{
										Name:      "testFunc",            // 8 bytes
										Arguments: "{\"key\":\"value\"}", // 15 bytes
									},
								},
							},
						},
					},
					{
						OfToolMessage: &adapter.ToolMessage{
							ToolCallID: "tool1", // 5 bytes
							Content: []adapter.ContentPartText{
								{
									Text: "Tool response", // 13 bytes
								},
							},
						},
					},
				},
				// No tokens set, will need estimation
			},
		},
	}

	// Total bytes in second frame: 13 + 29 + 16 + 5 + 8 + 15 + 5 + 13 = 104 bytes
	// With 2 bytes per token: 104 / 2 = 52 tokens
	// Expected total: 10 (first frame) + 52 (estimated) = 62 tokens
	expectedTokens := int64(62)

	tokens := thread.Tokens(samples)
	assert.Equal(t, expectedTokens, tokens, "Token count mismatch for multiple message types")
}

func TestThread_TokensEmptyFrames(t *testing.T) {
	// Create a test samples object with a known BytesPerTok value
	samples := client.NewSamples(5, 3.0)

	// Test with no frame containing token information
	thread := &Thread{
		Frames: []MessagesFrame{
			{
				Messages: []adapter.Message{
					{
						OfSystemMessage: &adapter.SystemMessage{
							Content: "System message", // 14 bytes
						},
					},
				},
				// No token count set
			},
			{
				Messages: []adapter.Message{
					{
						OfUserMessage: &adapter.UserMessage{
							Parts: []adapter.ContentPart{
								{
									OfContentPartText: &adapter.ContentPartText{
										Text: "User message", // 12 bytes
									},
								},
							},
						},
					},
				},
				// No token count set
			},
		},
	}

	// Total bytes: 14 + 12 = 26
	// With 3 bytes per token: 26 / 3 = ~8.67, which truncates to 8 tokens
	expectedTokens := int64(8)

	tokens := thread.Tokens(samples)
	assert.Equal(t, expectedTokens, tokens, "Token count mismatch for entirely estimated frames")
}

func TestThread_TokensWithRefusal(t *testing.T) {
	// Create a test samples object with a known BytesPerTok value
	samples := client.NewSamples(5, 2.5) // 2.5 bytes per token

	// Create a refusal message
	refusal := "I cannot comply with that request"

	// Create a thread with a refusal message
	thread := &Thread{
		Frames: []MessagesFrame{
			{
				Messages: []adapter.Message{
					{
						OfAssistantMessage: &adapter.AssistantMessage{
							Refusal: &refusal, // 33 bytes
						},
					},
				},
			},
		},
	}

	// Expected tokens: 33 / 2.5 = ~13.2, which rounds to 13 tokens when converted to int64
	expectedTokens := int64(13)

	tokens := thread.Tokens(samples)
	assert.Equal(t, expectedTokens, tokens, "Token count mismatch for refusal message")
}

func TestThread_TokensEmptyThread(t *testing.T) {
	// Create a test samples object
	samples := client.NewSamples(5, 3.0)

	// Test with empty frames list
	thread := &Thread{
		Frames: []MessagesFrame{},
	}

	// Expected tokens: 0
	expectedTokens := int64(0)

	tokens := thread.Tokens(samples)
	assert.Equal(t, expectedTokens, tokens, "Token count should be 0 for empty thread")
}
