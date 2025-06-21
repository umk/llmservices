package client

import (
	"fmt"
	"slices"
)

type Config struct {
	Preset *Preset `json:"preset" validate:"omitempty,min=1"`

	BaseURL string `json:"base_url" validate:"omitempty,url"`
	Key     string `json:"key" validate:"omitempty"`
	Model   string `json:"model" validate:"omitempty"`

	Concurrency int `json:"concurrency" validate:"omitempty,min=1"`
}

func getConfig(src *Config, allowed ...Preset) (*Config, error) {
	if src.Preset != nil && slices.Contains(allowed, *src.Preset) {
		return nil, fmt.Errorf("preset is not supported: %s", *src.Preset)
	}

	dest := Config{
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
	}

	if err := setConfig(&dest, src); err != nil {
		return nil, err
	}

	if dest.BaseURL == "" {
		dest.BaseURL = presetOpenAI.BaseURL
	}

	return &dest, nil
}

func setConfig(dest *Config, src *Config) error {
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
