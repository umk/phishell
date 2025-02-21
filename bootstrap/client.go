package bootstrap

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ClientRef struct {
	Config *ConfigService

	Client *openai.Client
}

type RequestCallback func(client *openai.Client) error

func NewClientRef(config *ConfigService) *ClientRef {
	ref := &ClientRef{
		Config: config,

		Client: getClient(config),
	}

	return ref
}

func (c *ClientRef) Request(ctx context.Context, cb RequestCallback) error {
	return cb(c.Client)
}

func getClient(config *ConfigService) *openai.Client {
	var opts []option.RequestOption

	if config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(config.BaseURL))
	}
	if config.Key != "" {
		opts = append(opts, option.WithAPIKey(config.Key))
	}

	return openai.NewClient(opts...)
}
