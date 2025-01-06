package bootstrap

import (
	"context"
	"errors"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/umk/phishell/util/termx"
)

type ClientRef struct {
	Config *ConfigService
	Client *openai.Client
}

type RequestCallback func(client *openai.Client) error

func NewClientRef(config *ConfigService) *ClientRef {
	ref := &ClientRef{Config: config}
	ref.refreshClient()

	return ref
}

func (c *ClientRef) Request(ctx context.Context, cb RequestCallback) error {
	for {
		if err := cb(c.Client); err == nil {
			return nil
		} else if err := c.clientRetryOrError(err); err != nil {
			return err
		}
	}
}

func (c *ClientRef) clientRetryOrError(err error) error {
	var apiErr *openai.Error

	if !(errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized) {
		return err
	}

	if c.Config.Key.Source != CfKeyChain {
		return err
	}

	currentKey, err := GetKey(c.Config.Profile)

	if err == nil && currentKey != c.Config.Key.Value {
		return nil
	}

	if err != nil || currentKey == c.Config.Key.Value {
		termx.Error.Println("The client could not be authorized.")

		s, err := ReadKeyAndUpdate(c.Config.Profile, false)
		if err != nil {
			return err
		}

		if s == "" {
			// Retry with the same key as before
			return nil
		}

		currentKey = s
	}

	c.Config.Key.Value = currentKey

	c.refreshClient()

	return nil
}

func (r *ClientRef) refreshClient() {
	var opts []option.RequestOption

	if r.Config.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(r.Config.BaseURL))
	}
	if r.Config.Key.Value != "" {
		opts = append(opts, option.WithAPIKey(r.Config.Key.Value))
	}

	r.Client = openai.NewClient(opts...)
}
