package host

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/umk/phishell/tool/host/process"
	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/util/execx"
)

type Provider struct {
	init sync.Mutex

	// Indicates whether provider termination was initiated, so no new
	// requests can be made to underlying process.
	terminated atomic.Bool

	// A prototype of the process' command.
	cmd *execx.Cmd

	process *process.Process

	// A pointer to provider info, which allows referencing provider
	// info by the host consumers without holding reference to provider
	// instance itself.
	info *ProviderInfo
}

type ProviderInfo struct {
	Pid           int
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

func (p *Provider) Info() *ProviderInfo {
	return p.info
}

func (p *Provider) Post(req *provider.Request) (*provider.Response, error) {
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
		log.Printf("process %d failed to terminate: %v\n", p.info.Pid, err)
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
		log.Printf("process %d timed out when terminating; sending SIGKILL\n", p.info.Pid)
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

	return p.info.Status
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

	p.info.Status = PsRunning

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
		return p.info.Status, false
	}

	p.info.Status = status
	p.info.StatusMessage = message

	log.Printf("process %d finalized with status %s: %s\n",
		p.info.Pid, p.info.Status, p.info.StatusMessage)

	p.process.Reset("provider terminated")

	return status, true
}
