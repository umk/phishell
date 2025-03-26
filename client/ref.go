package client

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/umk/phishell/config"
	"golang.org/x/sync/semaphore"
)

type Ref struct {
	Config  *config.Profile
	Client  openai.Client
	S       *semaphore.Weighted
	Samples *Samples
}

func NewRef(p *config.Profile) *Ref {
	return &Ref{
		Config:  p,
		Client:  getClient(p),
		S:       semaphore.NewWeighted(int64(p.Concurrency)),
		Samples: newSamples(samplesCount, defaultBytesPerTok),
	}
}

func getClient(config *config.Profile) openai.Client {
	var opts []option.RequestOption

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}
	if config.Key != "" {
		opts = append(opts, option.WithAPIKey(config.Key))
	}

	return openai.NewClient(opts...)
}
