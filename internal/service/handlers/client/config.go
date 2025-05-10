package client

import (
	"fmt"

	"github.com/umk/llmservices/pkg/client"
)

func getClientConfig(src *clientConfig) (*client.Config, error) {
	dest := client.Config{
		Concurrency: 1,
	}

	if src.Preset != nil {
		preset, ok := presets[*src.Preset]
		if !ok {
			return nil, fmt.Errorf("preset not found: %s", *src.Preset)
		}

		if err := setConfig(&dest, &preset); err != nil {
			return nil, err
		}
		dest.Preset = *src.Preset
	} else {
		dest.Preset = client.OpenAI
	}

	if err := setConfig(&dest, src); err != nil {
		return nil, err
	}

	if dest.BaseURL == "" {
		dest.BaseURL = presetOpenAI.BaseURL
	}

	return &dest, nil
}

func setConfig(dest *client.Config, src *clientConfig) error {
	if src.BaseURL != "" {
		dest.BaseURL = src.BaseURL
	}

	if src.Key != "" {
		dest.Key = src.Key
	}

	if src.Model != "" {
		dest.Model = src.Model
	}

	if src.Concurrency > 0 {
		dest.Concurrency = src.Concurrency
	}

	return nil
}
