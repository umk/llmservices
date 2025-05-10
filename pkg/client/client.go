package client

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/umk/llmservices/pkg/adapter"
	openaiadapter "github.com/umk/llmservices/pkg/adapter/openai"
	"golang.org/x/sync/semaphore"
)

type Config struct {
	Preset Preset

	BaseURL string
	Key     string
	Model   string

	Concurrency int
}

type Preset string

const (
	OpenAI Preset = "openai"
	Ollama Preset = "ollama"
)

type Client struct {
	config  *Config
	adapter adapter.Adapter
	s       *semaphore.Weighted
	Samples *Samples
}

const (
	defaultBytesPerTok = 3.25
	samplesCount       = 5
	minSampleSize      = 100
)

func New(p *Config) *Client {
	return &Client{
		config:  p,
		adapter: createAdapter(p),
		s:       semaphore.NewWeighted(int64(p.Concurrency)),
		Samples: newSamples(samplesCount, defaultBytesPerTok),
	}
}

func createAdapter(p *Config) adapter.Adapter {
	switch p.Preset {
	case OpenAI, Ollama:
		return createOpenAIAdapter(p)
	default:
		panic("preset is not supported")
	}
}

func createOpenAIAdapter(p *Config) adapter.Adapter {
	var opts []option.RequestOption

	if p.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(p.BaseURL))
	}
	if p.Key != "" {
		opts = append(opts, option.WithAPIKey(p.Key))
	}

	return &openaiadapter.Adapter{
		Client: openai.NewClient(opts...),
	}
}
