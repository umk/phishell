package host

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/tool/host/process"
)

type Provider struct {
	init sync.Mutex

	// Indicates whether provider termination was initiated, so no new
	// requests can be made to underlying process.
	terminated atomic.Bool

	// A pointer to provider info, which allows referencing provider
	// info by the host consumers without holding reference to provider
	// instance itself.
	Info *ProviderInfo

	process *process.Process
}

type ProviderInfo struct {
	Status        ProviderStatus
	StatusMessage string
}

type ProviderStatus int

const (
	_              ProviderStatus = iota
	PsInitializing                // Host application is reading headers
	PsRunning                     // Host application is reading messages
	PsCompleted                   // Process has exited normally
	PsFailed                      // Process was terminated or exited with non-zero code
)

func (s ProviderStatus) String() string {
	switch s {
	case PsInitializing:
		return "Initializing"
	case PsRunning:
		return "Running"
	case PsCompleted:
		return "Completed"
	case PsFailed:
		return "Failed"
	}

	return ""
}

func (p *Provider) Process() *process.Process {
	return p.process
}

func (p *Provider) Post(req *provider.ToolRequest) (*provider.ToolResponse, error) {
	if !p.terminated.Load() {
		return nil, errors.New("provider is terminated")
	}

	return p.process.Post(req)
}

func (p *Provider) Terminate(status ProviderStatus, message string) ProviderStatus {
	if s, ok := p.finalize(status, message); !ok {
		return s
	}

	c := p.process.Cmd()

	// Signal for graceful shutdown
	if err := c.Process.Signal(syscall.SIGTERM); err != nil {
		return status
	}

	// Wait for process completion
	done := make(chan error, 1)

	go func() {
		done <- p.process.WaitOnce()
		close(done)
	}()

	select {
	case <-done:
		// Do nothing
	case <-time.After(10 * time.Second):
		// If timeout, force terminate the process
		c.Process.Signal(syscall.SIGKILL)
	}

	return status
}

func (p *Provider) Wait() ProviderStatus {
	var status ProviderStatus
	var message string

	if err := p.process.WaitOnce(); err != nil {
		status = PsFailed

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			message = fmt.Sprintf("process exited with code %d", exitErr.ExitCode())
		} else {
			message = fmt.Sprintf("failed to complete process: %v", err)
		}
	} else {
		status = PsCompleted
		message = "process completed"
	}

	p.finalize(status, message)

	return p.Info.Status
}

func (p *Provider) initialize() error {
	p.init.Lock()
	defer p.init.Unlock()

	if p.terminated.Load() {
		return errors.New("provider is terminated")
	}

	if err := p.process.Init(); err != nil {
		return err
	}

	p.Info.Status = PsRunning

	go func() {
		if err := p.process.Read(); err != nil {
			p.Terminate(PsFailed, err.Error())
		}
	}()

	return nil
}

func (p *Provider) finalize(status ProviderStatus, message string) (s ProviderStatus, ok bool) {
	p.init.Lock()
	defer p.init.Unlock()

	if p.terminated.Swap(true) {
		// Just silently return as the process has already been finalized,
		// possibly with a different message and status.
		return p.Info.Status, false
	}

	p.Info.Status = status
	p.Info.StatusMessage = message

	p.process.Reset("provider terminated")

	return status, true
}
