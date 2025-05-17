package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umk/llmservices/internal/pointer"
	"github.com/umk/llmservices/pkg/client"
)

func TestGetClientConfig(t *testing.T) {
	tests := []struct {
		name      string
		input     *clientConfig
		expected  *client.Config
		expectErr bool
	}{
		{
			name:  "Default Configuration",
			input: &clientConfig{
				// No settings, should use default OpenAI preset
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     presetOpenAI.BaseURL,
				Concurrency: 1,
			},
			expectErr: false,
		},
		{
			name: "OpenAI Preset",
			input: &clientConfig{
				Preset: pointer.Ptr(client.OpenAI),
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     presetOpenAI.BaseURL,
				Concurrency: presetOpenAI.Concurrency,
			},
			expectErr: false,
		},
		{
			name: "Ollama Preset",
			input: &clientConfig{
				Preset: pointer.Ptr(client.Ollama),
			},
			expected: &client.Config{
				Preset:      client.Ollama,
				BaseURL:     presetOllama.BaseURL,
				Key:         presetOllama.Key,
				Concurrency: presetOllama.Concurrency,
			},
			expectErr: false,
		},
		{
			name: "Invalid Preset",
			input: &clientConfig{
				Preset: pointer.Ptr(client.Preset("invalid")),
			},
			expectErr: true,
		},
		{
			name: "Custom Configuration without Preset",
			input: &clientConfig{
				BaseURL:     "https://custom-api.example.com/v1/",
				Key:         "test-key",
				Model:       "test-model",
				Concurrency: 10,
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://custom-api.example.com/v1/",
				Key:         "test-key",
				Model:       "test-model",
				Concurrency: 10,
			},
			expectErr: false,
		},
		{
			name: "Preset with Overrides",
			input: &clientConfig{
				Preset:      pointer.Ptr(client.OpenAI),
				Key:         "custom-key",
				Model:       "custom-model",
				Concurrency: 3,
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     presetOpenAI.BaseURL,
				Key:         "custom-key",
				Model:       "custom-model",
				Concurrency: 3,
			},
			expectErr: false,
		},
		{
			name: "Empty BaseURL Falls Back to Default",
			input: &clientConfig{
				Key:         "test-key",
				Model:       "test-model",
				Concurrency: 2,
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     presetOpenAI.BaseURL,
				Key:         "test-key",
				Model:       "test-model",
				Concurrency: 2,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getClientConfig(tt.input)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.Preset, got.Preset, "Preset mismatch")
			assert.Equal(t, tt.expected.BaseURL, got.BaseURL, "BaseURL mismatch")
			assert.Equal(t, tt.expected.Key, got.Key, "Key mismatch")
			assert.Equal(t, tt.expected.Model, got.Model, "Model mismatch")
			assert.Equal(t, tt.expected.Concurrency, got.Concurrency, "Concurrency mismatch")
		})
	}
}

func TestSetConfig(t *testing.T) {
	tests := []struct {
		name     string
		dest     *client.Config
		src      *clientConfig
		expected *client.Config
	}{
		{
			name: "Empty Source",
			dest: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 1,
			},
			src: &clientConfig{},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 1,
			},
		},
		{
			name: "Full Override",
			dest: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 1,
			},
			src: &clientConfig{
				BaseURL:     "https://new.example.com/",
				Key:         "new-key",
				Model:       "new-model",
				Concurrency: 5,
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://new.example.com/",
				Key:         "new-key",
				Model:       "new-model",
				Concurrency: 5,
			},
		},
		{
			name: "Partial Override",
			dest: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 1,
			},
			src: &clientConfig{
				Key:   "new-key",
				Model: "new-model",
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "new-key",
				Model:       "new-model",
				Concurrency: 1,
			},
		},
		{
			name: "Only Concurrency Change",
			dest: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 1,
			},
			src: &clientConfig{
				Concurrency: 10,
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 10,
			},
		},
		{
			name: "Invalid Concurrency (Zero)",
			dest: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 5,
			},
			src: &clientConfig{
				Concurrency: 0,
			},
			expected: &client.Config{
				Preset:      client.OpenAI,
				BaseURL:     "https://original.example.com/",
				Key:         "original-key",
				Model:       "original-model",
				Concurrency: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone the destination to avoid modifying the test case
			dest := *tt.dest

			err := setConfig(&dest, tt.src)

			// setConfig should always succeed as currently implemented
			require.NoError(t, err)

			assert.Equal(t, tt.expected.BaseURL, dest.BaseURL, "BaseURL mismatch")
			assert.Equal(t, tt.expected.Key, dest.Key, "Key mismatch")
			assert.Equal(t, tt.expected.Model, dest.Model, "Model mismatch")
			assert.Equal(t, tt.expected.Concurrency, dest.Concurrency, "Concurrency mismatch")
			assert.Equal(t, tt.expected.Preset, dest.Preset, "Preset should never change")
		})
	}
}
