package host

import (
	"context"

	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/tool/host/provider/schema"
)

type HostedToolHandler struct {
	provider *provider.Provider
	req      *schema.Request
}

func (h *HostedToolHandler) Execute(ctx context.Context) (any, error) {
	res, err := h.provider.Post(h.req)
	if err != nil {
		return nil, err
	}

	return res.Content, nil
}
