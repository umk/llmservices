package openai

import (
	"github.com/openai/openai-go"
	"github.com/umk/llmservices/pkg/adapter"
)

func getContentPart(part *adapter.ContentPart) openai.ChatCompletionContentPartUnionParam {
	switch {
	case part.OfContentPartText != nil:
		return getTextContentPart(part)
	case part.OfContentPartImageUrl != nil:
		return getImageContentPart(part)
	default:
		return openai.ChatCompletionContentPartUnionParam{}
	}
}

func getTextContentPart(part *adapter.ContentPart) openai.ChatCompletionContentPartUnionParam {
	return openai.ChatCompletionContentPartUnionParam{
		OfText: &openai.ChatCompletionContentPartTextParam{
			Text: part.OfContentPartText.Text,
		},
	}
}

func getImageContentPart(part *adapter.ContentPart) openai.ChatCompletionContentPartUnionParam {
	return openai.ChatCompletionContentPartUnionParam{
		OfImageURL: &openai.ChatCompletionContentPartImageParam{
			ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
				URL: part.OfContentPartImageUrl.ImageUrl,
			},
		},
	}
}
