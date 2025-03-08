package host

import (
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/tool/builtin"
	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/tool/host/provider/schema"
)

func (h *Host) Tools() ([]openai.ChatCompletionToolParam, error) {
	tools, err := h.getProviderTools()
	if err != nil {
		return nil, err
	}

	r := make([]openai.ChatCompletionToolParam, 0, len(tools))

	// The tools map is immutable and can be accessed concurrently
	for _, tool := range tools {
		r = append(r, tool.param)
	}

	r = append(r, builtin.Tools...)

	return r, nil
}

func (h *Host) Get(f *openai.ChatCompletionMessageToolCallFunction) (tool.Handler, error) {
	tools, err := h.getProviderTools()
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot get working directory: %w", err)
	}

	// Built-in commands
	switch f.Name {
	case builtin.FsCreateOrUpdateToolName:
		return builtin.NewFsCreateOrUpdateToolHandler(f.Arguments, wd)
	case builtin.FsReadToolName:
		return builtin.NewFsReadToolHandler(f.Arguments, wd)
	case builtin.FsDeleteToolName:
		return builtin.NewFsDeleteToolHandler(f.Arguments, wd)
	case builtin.ExecCommandToolName:
		return builtin.NewExecCommandToolHandler(f.Arguments, wd)
	}

	// The tools map is immutable and can be accessed concurrently
	t, ok := tools[f.Name]
	if !ok {
		return nil, fmt.Errorf("no handler registered for %s", f.Name)
	}

	req := &schema.Request{
		CallID: uuid.NewString(),
		Function: schema.Function{
			Name:      f.Name,
			Arguments: f.Arguments,
		},
		Context: schema.Context{
			Dir: wd,
		},
	}

	return &HostedToolHandler{
		provider: t.provider,
		req:      req,
	}, nil
}

func (h *Host) providerAdd(p *provider.Provider) {
	h.providers = append(h.providers, p)

	h.refreshTools()

	h.wg.Add(1)
}

func (h *Host) providerDel(p *provider.Provider) {
	h.providers = slices.DeleteFunc(h.providers, func(current *provider.Provider) bool {
		return current == p
	})

	h.refreshTools()

	h.wg.Done()
}

func (h *Host) getProviderTools() (toolsMap, error) {
	if err, ok := h.tools.(error); ok {
		return nil, err
	}

	return h.tools.(toolsMap), nil
}

func (f *Host) refreshTools() {
	tools := make(map[string]*providerTool)

	for _, p := range f.providers {
		if p.Info.Status != provider.PsRunning {
			continue
		}

		for k, v := range p.Process.Tools() {
			if _, ok := tools[k]; ok {
				f.tools = fmt.Errorf("duplicate exports of the tool: %s", k)
				return
			} else {
				tools[k] = &providerTool{
					provider: p,
					param:    v,
				}
			}
		}
	}

	f.tools = tools
}
