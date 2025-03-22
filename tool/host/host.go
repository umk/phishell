package host

import (
	"errors"
	"os/exec"
	"sync"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/tool/host/provider"
)

type Host struct {
	mu sync.Mutex

	// Indicates whether host termination was initiated, so no new
	// providers can be executed.
	terminated bool

	providers []*provider.Provider
	tools     any // either tools map or error

	// Counts running providers to wait before close the host.
	wg sync.WaitGroup
}

type toolsMap = map[string]*providerTool

type providerTool struct {
	provider *provider.Provider
	param    openai.ChatCompletionToolParam
}

func NewHost() *Host {
	return &Host{tools: make(map[string]*providerTool)}
}

func (h *Host) Execute(cmd *exec.Cmd) (*provider.Provider, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.terminated {
		return nil, errors.New("host is terminated")
	}

	p, err := provider.Start(cmd)
	if err != nil {
		return nil, err
	}

	h.providerAdd(p)

	go func() {
		p.Wait()

		h.mu.Lock()
		defer h.mu.Unlock()

		h.providerDel(p)
	}()

	return p, nil
}

func (h *Host) Close() error {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		h.terminated = true

		for _, p := range h.providers {
			go p.Terminate(provider.PsCompleted, "host terminated")
		}
	}()

	h.wg.Wait()

	return nil
}
