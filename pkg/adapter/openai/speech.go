package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/openai/openai-go"
	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Adapter) Speech(ctx context.Context, message adapter.SpeechMessage, params adapter.SpeechParams) (adapter.Speech, error) {
	resp, err := c.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Input:          message.Content,
		Instructions:   getOpt(message.Instructions),
		Model:          params.Model,
		Voice:          openai.AudioSpeechNewParamsVoice(params.Voice),
		Speed:          getOpt(params.Speed),
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatPCM,
	})
	if err != nil {
		return adapter.Speech{}, err
	}

	n := config.C.AudioBufSize
	if resp.ContentLength != -1 {
		// Ensure content length is not greater than 1GB
		if resp.ContentLength > 1<<30 {
			return adapter.Speech{}, fmt.Errorf("response too large: %d bytes", resp.ContentLength)
		}
		n = int(resp.ContentLength)
	}

	buf := bytes.NewBuffer(make([]byte, 0, n))

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return adapter.Speech{}, err
	}

	return adapter.Speech{
		Audio: buf.Bytes(),
	}, nil
}
