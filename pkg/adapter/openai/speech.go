package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"

	"github.com/openai/openai-go"
	"github.com/umk/llmservices/internal/audio"
	"github.com/umk/llmservices/internal/config"
	"github.com/umk/llmservices/pkg/adapter"
)

func (c *Adapter) Speech(ctx context.Context, message adapter.SpeechMessage, params adapter.SpeechParams) (adapter.Speech, error) {
	var speed *float64
	if params.Speed != nil {
		s := math.Max(0.25, math.Min(4.0, *params.Speed))
		speed = &s
	}

	resp, err := c.Audio.Speech.New(ctx, openai.AudioSpeechNewParams{
		Input:          message.Content,
		Instructions:   getOpt(message.Instructions),
		Model:          params.Model,
		Voice:          openai.AudioSpeechNewParamsVoice(params.Voice),
		Speed:          getOpt(speed),
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatWAV,
	})
	if err != nil {
		return adapter.Speech{}, err
	}

	n := config.Cur.AudioBufSize
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

	a, err := audio.ParseWAV(buf.Bytes())
	if err != nil {
		return adapter.Speech{}, err
	}

	return adapter.Speech{Audio: a}, nil
}
