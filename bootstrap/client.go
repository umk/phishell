package bootstrap

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Profile struct {
	Config *ConfigProfile

	Client *openai.Client
}

type RequestCallback func(client *openai.Client) error

func NewProfile(config *ConfigProfile) *Profile {
	return &Profile{
		Config: config,

		Client: getClient(config),
	}
}

func (c *Profile) Request(ctx context.Context, cb RequestCallback) error {
	return cb(c.Client)
}

func getClient(profile *ConfigProfile) *openai.Client {
	var opts []option.RequestOption

	if profile.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(profile.BaseURL))
	}
	if profile.Key != "" {
		opts = append(opts, option.WithAPIKey(profile.Key))
	}

	return openai.NewClient(opts...)
}
