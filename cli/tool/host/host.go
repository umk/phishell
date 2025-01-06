package host

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/util/execx"
)

type Host struct {
	mu sync.Mutex

	processes []*ToolProcess
	tools     map[string]*processTool

	wg sync.WaitGroup
}

type processTool struct {
	Process *ToolProcess
	Tool    openai.ChatCompletionToolParam
}

func NewHost() *Host {
	return &Host{tools: make(map[string]*processTool)}
}

func (h *Host) Execute(ctx context.Context, c *execx.Cmd) (*ToolProcess, error) {
	cmd := c.Command()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe Stdin: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe Stdout: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	p := &ToolProcess{
		host: h,

		cmd:   cmd,
		stdin: stdin,

		requests: make(map[string]chan<- any),

		Tools: make(map[string]openai.ChatCompletionToolParam),
		Info: &ToolProcessInfo{
			Status: TsInitializing,
		},
	}

	if err := p.initialize(ctx, stdout); err != nil {
		p.Terminate(TsFailed, err.Error())
		return nil, err
	}

	h.add(p)

	h.wg.Add(1)

	go func() {
		p.Wait()
		h.wg.Done()
	}()

	return p, nil
}

func (h *Host) Close() error {
	go func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		for _, process := range h.processes {
			if process.Info.Status == TsRunning {
				go process.Terminate(TsCompleted, "")
			}
		}
	}()

	h.wg.Wait()

	return nil
}
