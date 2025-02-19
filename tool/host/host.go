package host

import (
	"errors"
	"sync"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/tool/host/process"
	"github.com/umk/phishell/util/execx"
)

type Host struct {
	mu sync.Mutex

	// Indicates whether host termination was initiated, so no new
	// providers can be executed.
	terminated bool

	providers []*Provider
	tools     any // either tools map or error

	// Counts running providers to wait before close the host.
	wg sync.WaitGroup
}

type toolsMap = map[string]*providerTool

type providerTool struct {
	provider *Provider
	param    openai.ChatCompletionToolParam
}

func NewHost() *Host {
	return &Host{tools: make(map[string]*providerTool)}
}

func (h *Host) Execute(c *execx.Cmd) (*Provider, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.terminated {
		return nil, errors.New("host is terminated")
	}

	// Start provider process and create provider
	pr, err := process.Start(c)
	if err != nil {
		return nil, err
	}

	p := &Provider{
		cmd:     c,
		process: pr,

		info: &ProviderInfo{
			Pid:    pr.Cmd().Process.Pid,
			Status: PsInitializing,
		},
	}

	// Register created provider in the host
	h.providerAdd(p)

	go func() {
		p.Wait()

		h.mu.Lock()
		defer h.mu.Unlock()

		h.providerDel(p)
	}()

	if err := p.initialize(); err != nil {
		p.Terminate(PsFailed, err.Error())
		return nil, err
	}

	return p, nil
}

func (h *Host) Close() error {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		h.terminated = true

		for _, p := range h.providers {
			go p.Terminate(PsCompleted, "host terminated")
		}
	}()

	h.wg.Wait()

	return nil
}
