package host

import (
	"context"

	"github.com/umk/phishell/tool/host/provider"
)

type HostedToolHandler struct {
	provider *Provider
	req      *provider.Request
}

func (h *HostedToolHandler) Execute(ctx context.Context) (any, error) {
	res, err := h.provider.process.Post(h.req)
	if err != nil {
		return nil, err
	}

	return res.Content, nil
}
