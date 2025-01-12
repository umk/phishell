package host

import (
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/tool"
	"github.com/umk/phishell/tool/builtin"
)

func (h *Host) Tools() []openai.ChatCompletionToolParam {
	var r []openai.ChatCompletionToolParam

	for _, tool := range h.tools {
		r = append(r, tool.Tool)
	}

	r = append(r, builtin.Tools...)

	return r
}

func (h *Host) Get(f *openai.ChatCompletionMessageToolCallFunction) (tool.Handler, error) {
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
	case builtin.ExecHttpCallToolName:
		return builtin.NewExecHttpCallToolHandler(f.Arguments)
	}

	// The tools map is immutable and can be accessed concurrently
	t, ok := h.tools[f.Name]
	if !ok {
		return nil, fmt.Errorf("no handler registered for %s", f.Name)
	}

	return &HostedToolHandler{
		process: t.Process,
		req: &provider.ToolRequest{
			CallID: uuid.NewString(),
			Function: provider.ToolRequestFunction{
				Name:      f.Name,
				Arguments: f.Arguments,
			},
			Context: provider.ToolRequestContext{
				Dir: wd,
			},
		},
	}, nil
}

func (h *Host) add(p *ToolProcess) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.processes = append(h.processes, p)

	h.refreshTools()
}

func (h *Host) delete(p *ToolProcess) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.processes = slices.DeleteFunc(h.processes, func(current *ToolProcess) bool {
		return current == p
	})

	h.refreshTools()
}

func (f *Host) refreshTools() {
	tools := make(map[string]*processTool)

	for _, p := range f.processes {
		if p.Info.Status != TsRunning {
			continue
		}

		for k, v := range p.Tools {
			tools[k] = &processTool{
				Process: p,
				Tool:    v,
			}
		}
	}

	f.tools = tools
}
