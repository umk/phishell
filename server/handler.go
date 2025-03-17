package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/umk/phishell/client"
)

type handler struct {
	clients map[string]*Client
}

func newHandler() (*handler, error) {
	h := &handler{
		clients: make(map[string]*Client),
	}
	for id, p := range client.Profiles {
		c, err := NewClient(p.Config.BaseURL, WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
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

func (h *handler) CreateChatCompletion(w http.ResponseWriter, r *http.Request) {
	req, ok := parseRequestBody[CreateChatCompletionRequest](w, r)
	if !ok {
		return
	}

	p, _ := req.Model.AsCreateChatCompletionRequestModel0()
	base, ok := h.clients[p]
	if !ok {
		http.Error(w, "Invalid profile: "+p, http.StatusBadRequest)
		return
	}

	c := client.Get(client.Profiles[p])
	if err := c.S.Acquire(r.Context(), 1); err != nil {
		http.Error(w, "Error acquiring semaphore: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.S.Release(1)

	res, err := base.CreateChatCompletion(r.Context(), req)
	if err != nil {
		http.Error(w, "Error creating chat completion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	propagateResponse(w, res)
}

func (h *handler) CreateCompletion(w http.ResponseWriter, r *http.Request) {
	req, ok := parseRequestBody[CreateCompletionRequest](w, r)
	if !ok {
		return
	}

	p, _ := req.Model.AsCreateCompletionRequestModel0()
	base, ok := h.clients[p]
	if !ok {
		http.Error(w, "Invalid profile: "+p, http.StatusBadRequest)
		return
	}

	c := client.Get(client.Profiles[p])
	if err := c.S.Acquire(r.Context(), 1); err != nil {
		http.Error(w, "Error acquiring semaphore: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.S.Release(1)

	res, err := base.CreateCompletion(r.Context(), req)
	if err != nil {
		http.Error(w, "Error creating completion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	propagateResponse(w, res)
}

func (h *handler) CreateEmbedding(w http.ResponseWriter, r *http.Request) {
	req, ok := parseRequestBody[CreateEmbeddingRequest](w, r)
	if !ok {
		return
	}

	p, _ := req.Model.AsCreateEmbeddingRequestModel0()
	base, ok := h.clients[p]
	if !ok {
		http.Error(w, "Invalid profile: "+p, http.StatusBadRequest)
		return
	}

	c := client.Get(client.Profiles[p])
	if err := c.S.Acquire(r.Context(), 1); err != nil {
		http.Error(w, "Error acquiring semaphore: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer c.S.Release(1)

	res, err := base.CreateEmbedding(r.Context(), req)
	if err != nil {
		http.Error(w, "Error creating embedding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	propagateResponse(w, res)
}
