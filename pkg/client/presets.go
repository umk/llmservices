package client

type Preset string

const (
	OpenAI Preset = "openai"
	Ollama Preset = "ollama"
)

var presetOpenAI = Config{
	BaseURL:     "https://api.openai.com/v1/",
	Concurrency: 5,
}

var presetOllama = Config{
	BaseURL:     "http://localhost:11434/v1/",
	Key:         "ollama",
	Concurrency: 1,
}

var presets = map[Preset]Config{
	OpenAI: presetOpenAI,
	Ollama: presetOllama,
}
