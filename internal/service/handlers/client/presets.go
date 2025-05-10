package client

import "github.com/umk/llmservices/pkg/client"

var presetOpenAI = clientConfig{
	BaseURL:     "https://api.openai.com/v1/",
	Concurrency: 5,
}

var presetOllama = clientConfig{
	BaseURL:     "http://localhost:11434/v1/",
	Key:         "ollama",
	Concurrency: 1,
}

var presets = map[client.Preset]clientConfig{
	client.OpenAI: presetOpenAI,
	client.Ollama: presetOllama,
}
