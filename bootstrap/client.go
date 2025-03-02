package bootstrap

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var Clients []*ClientRef

func InitClients() {
	for _, p := range Config.Services {
		Clients = append(Clients, &ClientRef{
			Config: p,
			Client: getClient(p),
		})
	}
}

type ClientRef struct {
	Config *ConfigService
	Client *openai.Client
}

type RequestCallback func(client *openai.Client) error

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

// GetClient gets the default client to use outside of the chat context
// where user can pick the client explicitly.
func GetDefaultClient() *ClientRef {
	if len(Clients) == 0 {
		panic("no clients defined")
	}

	return Clients[0]
}
