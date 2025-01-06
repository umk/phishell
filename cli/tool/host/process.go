package host

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/openai/openai-go"
)

type ToolProcess struct {
	mu   sync.Mutex
	host *Host

	cmd   *exec.Cmd
	stdin io.WriteCloser

	requests map[string]chan<- any

	Tools map[string]openai.ChatCompletionToolParam

	Info *ToolProcessInfo
}

type ToolProcessInfo struct {
	Status        ToolStatus
	StatusMessage string
}

type ToolStatus int

const (
	_              ToolStatus = iota
	TsInitializing            // Host application is reading headers
	TsRunning                 // Host application is reading messages
	TsCompleted               // Process has exited normally
	TsFailed                  // Process was terminated or exited with non-zero code
)

func (s ToolStatus) String() string {
	switch s {
	case TsInitializing:
		return "Initializing"
	case TsRunning:
		return "Running"
	case TsCompleted:
		return "Completed"
	case TsFailed:
		return "Failed"
	}

	return ""
}

func (p *ToolProcess) initialize(ctx context.Context, stdout io.Reader) error {
	scanner := bufio.NewScanner(stdout)

	init := make(chan error)

	go func() {
		init <- p.readHeader(scanner)
	}()

	select {
	case err := <-init:
		if err != nil {
			return fmt.Errorf("failed to initialize tools: %w", err)
		}
		p.Info.Status = TsRunning
	case <-time.After(10 * time.Second):
		return errors.New("initialization timeout")
	}

	go func() {
		if err := p.readMessages(ctx, scanner); err != nil {
			p.Terminate(TsFailed, err.Error())
		}
	}()

	return nil
}

func (p *ToolProcess) Terminate(status ToolStatus, message string) {
	p.finalize(status, message)

	p.cmd.Process.Kill()
}

func (p *ToolProcess) Wait() ToolStatus {
	var status ToolStatus
	var message string

	if err := p.cmd.Wait(); err != nil {
		status = TsFailed

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			message = fmt.Sprintf("process exited with code %d", exitErr.ExitCode())
		} else {
			message = fmt.Sprintf("failed to complete process: %s", err)
		}
	} else {
		status = TsCompleted
	}

	p.finalize(status, message)

	return p.Info.Status
}

func (p *ToolProcess) finalize(status ToolStatus, message string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Info.Status != TsInitializing && p.Info.Status != TsRunning {
		// Just silently return as the process has already been finalized,
		// possibly with a different message and status.
		return
	}

	p.Info.Status = status
	p.Info.StatusMessage = message

	p.host.delete(p)

	// Close pending requests
	for _, ch := range p.requests {
		ch <- errors.New("handler is terminated")
	}

	p.requests = make(map[string]chan<- any)
}

func (p *ToolProcess) Pid() int {
	return p.cmd.Process.Pid
}
