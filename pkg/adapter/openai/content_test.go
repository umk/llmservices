package openai

import (
	"testing"

	"github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
	"github.com/umk/llmservices/pkg/adapter"
)

func TestGetContentPart(t *testing.T) {
	tests := []struct {
		name     string
		input    *adapter.ContentPart
		expected openai.ChatCompletionContentPartUnionParam
	}{
		{
			name: "Text content part",
			input: &adapter.ContentPart{
				OfContentPartText: &adapter.ContentPartText{
					Text: "Hello, world!",
				},
			},
			expected: openai.ChatCompletionContentPartUnionParam{
				OfText: &openai.ChatCompletionContentPartTextParam{
					Text: "Hello, world!",
				},
			},
		},
		{
			name: "Image URL content part",
			input: &adapter.ContentPart{
				OfContentPartImageUrl: &adapter.ContentPartImage{
					ImageUrl: "https://example.com/image.jpg",
				},
			},
			expected: openai.ChatCompletionContentPartUnionParam{
				OfImageURL: &openai.ChatCompletionContentPartImageParam{
					ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
						URL: "https://example.com/image.jpg",
					},
				},
			},
		},
		{
			name:     "Empty content part",
			input:    &adapter.ContentPart{},
			expected: openai.ChatCompletionContentPartUnionParam{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getContentPart(tc.input)

			if tc.expected.OfText != nil {
				// Check if Text content part is as expected
				assert.NotNil(t, result.OfText)
				assert.Equal(t, tc.expected.OfText.Text, result.OfText.Text)
				assert.Nil(t, result.OfImageURL)
			} else if tc.expected.OfImageURL != nil {
				// Check if Image URL content part is as expected
				assert.NotNil(t, result.OfImageURL)
				assert.Equal(t, tc.expected.OfImageURL.ImageURL.URL, result.OfImageURL.ImageURL.URL)
				assert.Nil(t, result.OfText)
			} else {
				// Check if it's empty as expected
				assert.Nil(t, result.OfText)
				assert.Nil(t, result.OfImageURL)
			}
		})
	}
}

func TestGetTextContentPart(t *testing.T) {
	part := &adapter.ContentPart{
		OfContentPartText: &adapter.ContentPartText{
			Text: "Sample text",
		},
	}

	result := getTextContentPart(part)

	assert.NotNil(t, result.OfText)
	assert.Equal(t, "Sample text", result.OfText.Text)
	assert.Nil(t, result.OfImageURL)
}

func TestGetImageContentPart(t *testing.T) {
	part := &adapter.ContentPart{
		OfContentPartImageUrl: &adapter.ContentPartImage{
			ImageUrl: "https://example.com/sample.png",
		},
	}

	result := getImageContentPart(part)

	assert.NotNil(t, result.OfImageURL)
	assert.Equal(t, "https://example.com/sample.png", result.OfImageURL.ImageURL.URL)
	assert.Nil(t, result.OfText)
}
