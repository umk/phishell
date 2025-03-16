package client

import (
	"context"
	"errors"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/umk/phishell/config"
)

var ChatProfiles []*Ref
var Profiles = make(map[string]*Ref)

var Default *Ref

func Init() error {
	if len(config.Config.ChatProfiles) == 0 {
		return errors.New("no chat profiles defined")
	}

	for _, p := range config.Config.Profiles {
		Profiles[p.Profile] = &Ref{
			Config: p,
			Client: getClient(p),
		}
	}

	for _, id := range config.Config.ChatProfiles {
		ChatProfiles = append(ChatProfiles, Profiles[id])
	}

	Default = ChatProfiles[0]
	return nil
}

type Ref struct {
	Config *config.Profile
	Client *openai.Client
}

type RequestCallback func(client *openai.Client) error

func (c *Ref) Request(ctx context.Context, cb RequestCallback) error {
	return cb(c.Client)
}

func getClient(config *config.Profile) *openai.Client {
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
func GetDefaultClient() *Ref {
	if len(ChatProfiles) == 0 {
		panic("no clients defined")
	}

	return ChatProfiles[0]
}
