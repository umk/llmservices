package client

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/umk/llmservices/pkg/adapter"
	openaiadapter "github.com/umk/llmservices/pkg/adapter/openai"
	"golang.org/x/sync/semaphore"
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

func New(p *Config) (*Client, error) {
	a, err := GetAdapter(p)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:  p,
		adapter: a,
		s:       semaphore.NewWeighted(int64(p.Concurrency)),
		Samples: NewSamples(samplesCount, defaultBytesPerTok),
	}, nil
}

func GetAdapter(p *Config) (adapter.Adapter, error) {
	if p.Preset == nil {
		return GetOpenAIAdapter(p)
	}

	switch *p.Preset {
	case OpenAI, Ollama:
		return GetOpenAIAdapter(p)
	default:
		panic("preset is not supported")
	}
}

func GetOpenAIAdapter(p *Config) (adapter.Adapter, error) {
	p, err := getConfig(p, OpenAI, Ollama)
	if err != nil {
		return nil, err
	}

	var opts []option.RequestOption

	if p.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(p.BaseURL))
	}
	if p.Key != "" {
		opts = append(opts, option.WithAPIKey(p.Key))
	}

	return &openaiadapter.Adapter{
		Client: openai.NewClient(opts...),
	}, nil
}
