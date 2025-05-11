package client

import (
	"testing"

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

			// Check error status
			if (err != nil) != tt.expectErr {
				t.Errorf("getClientConfig() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if tt.expectErr {
				return
			}

			// Check returned config matches expected
			if got.Preset != tt.expected.Preset {
				t.Errorf("Preset = %v, want %v", got.Preset, tt.expected.Preset)
			}

			if got.BaseURL != tt.expected.BaseURL {
				t.Errorf("BaseURL = %v, want %v", got.BaseURL, tt.expected.BaseURL)
			}

			if got.Key != tt.expected.Key {
				t.Errorf("Key = %v, want %v", got.Key, tt.expected.Key)
			}

			if got.Model != tt.expected.Model {
				t.Errorf("Model = %v, want %v", got.Model, tt.expected.Model)
			}

			if got.Concurrency != tt.expected.Concurrency {
				t.Errorf("Concurrency = %v, want %v", got.Concurrency, tt.expected.Concurrency)
			}
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
			if err != nil {
				t.Errorf("setConfig() error = %v", err)
				return
			}

			// Check if fields were correctly updated
			if dest.BaseURL != tt.expected.BaseURL {
				t.Errorf("BaseURL = %v, want %v", dest.BaseURL, tt.expected.BaseURL)
			}

			if dest.Key != tt.expected.Key {
				t.Errorf("Key = %v, want %v", dest.Key, tt.expected.Key)
			}

			if dest.Model != tt.expected.Model {
				t.Errorf("Model = %v, want %v", dest.Model, tt.expected.Model)
			}

			if dest.Concurrency != tt.expected.Concurrency {
				t.Errorf("Concurrency = %v, want %v", dest.Concurrency, tt.expected.Concurrency)
			}

			// Preset should never change in setConfig
			if dest.Preset != tt.expected.Preset {
				t.Errorf("Preset was modified: got %v, want %v", dest.Preset, tt.expected.Preset)
			}
		})
	}
}
