package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/umk/phishell/client"
	"github.com/umk/phishell/server/internal"
	"github.com/umk/phishell/util/stringsx"
	"github.com/umk/phishell/util/termx"
)

type proxyHandler struct {
	clients map[string]*internal.Client
}

func newProxyHandler() (*proxyHandler, error) {
	h := &proxyHandler{
		clients: make(map[string]*internal.Client),
	}
	for id, p := range client.Profiles {
		c, err := internal.NewClient(p.Config.BaseURL, internal.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", "Bearer "+p.Config.Key)
			return nil
		}))
		if err != nil {
			return nil, fmt.Errorf("error creating client for %s: %w", id, err)
		}

		h.clients[id] = c
	}

	return h, nil
}

func (h *proxyHandler) CreateChatCompletion(w http.ResponseWriter, r *http.Request) {
	req, ok := parseRequestBody[internal.CreateChatCompletionRequest](w, r)
	if !ok {
		return
	}

	base, ok := h.clients[req.V.Model]
	if !ok {
		http.Error(w, "Invalid profile: "+req.V.Model, http.StatusBadRequest)
		return
	}

	c := client.Get(client.Profiles[req.V.Model])
	if err := c.S.Acquire(r.Context(), 1); err != nil {
		http.Error(w, "Error creating chat completion: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.S.Release(1)

	res, err := base.CreateChatCompletion(r.Context(), req)
	if err != nil {
		http.Error(w, "Error creating chat completion: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	v, ok := parseResponseBody[internal.CreateChatCompletionResponse](w, res)
	if !ok {
		return
	}

	if len(v.V.Choices) == 1 {
		c := v.V.Choices[0]
		if tcs := c.V.Message.V.ToolCalls; tcs != nil {
			if len(*tcs) > 1 {
				// Can't handle approve and revise loop of multiple calls, so
				// force the number of calls to one.
				*tcs = (*tcs)[:1]
			}

			if len(*tcs) == 1 {
				tc := (*tcs)[0]
				toks := stringsx.Tokens(tc.V.Function.V.Name)
				termx.MD.Printf("Running: %s", stringsx.DisplayName(toks))
			}
		}
	}

	propagateResponse(w, res, v)
}

func (h *proxyHandler) CreateEmbedding(w http.ResponseWriter, r *http.Request) {
	req, ok := parseRequestBody[internal.CreateEmbeddingRequest](w, r)
	if !ok {
		return
	}

	base, ok := h.clients[req.V.Model]
	if !ok {
		http.Error(w, "Invalid profile: "+req.V.Model, http.StatusBadRequest)
		return
	}

	c := client.Get(client.Profiles[req.V.Model])
	if err := c.S.Acquire(r.Context(), 1); err != nil {
		http.Error(w, "Error creating embedding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.S.Release(1)

	res, err := base.CreateEmbedding(r.Context(), req)
	if err != nil {
		http.Error(w, "Error creating embedding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	propagateResponse(w, res, res.Body)
}
