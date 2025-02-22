package provider

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
	"github.com/umk/phishell/tool/host/provider/schema"
	"github.com/umk/phishell/util/execx"
)

const RestartExitCode = 42

type Provider struct {
	init sync.Mutex

	// Indicates whether provider termination was initiated, so no new
	// requests can be made to underlying process.
	terminated atomic.Bool

	// A prototype of the process' command.
	Cmd *execx.Cmd

	Process *process.Process

	// A pointer to provider info, which allows referencing provider
	// info by the host consumers without holding reference to provider
	// instance itself.
	Info *Info
}

type Info struct {
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

func (p *Provider) Post(req *schema.Request) (*schema.Response, error) {
	if !p.terminated.Load() {
		return nil, errors.New("provider is terminated")
	}

	return p.Process.Post(req)
}

func (p *Provider) Terminate(status ProviderStatus, message string) ProviderStatus {
	if s, ok := p.finalize(status, message); !ok {
		return s
	}

	c := p.Process.Cmd()

	// Signal for graceful shutdown
	if err := c.Process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("process %d failed to terminate: %v\n", p.Info.Pid, err)
		return status
	}

	// Wait for process completion
	done := make(chan error, 1)

	go func() {
		done <- p.Process.WaitOnce()
		close(done)
	}()

	select {
	case <-done:
		// Do nothing
	case <-time.After(10 * time.Second):
		// If timeout, force terminate the process
		log.Printf("process %d timed out when terminating; sending SIGKILL\n", p.Info.Pid)
		c.Process.Signal(syscall.SIGKILL)
	}

	return status
}

func (p *Provider) Wait() ProviderStatus {
	var status ProviderStatus
	var message string

	for {
		if err := p.Process.WaitOnce(); err != nil {
			status = PsFailed

			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				if exitErr.ExitCode() == RestartExitCode {
					if err := p.restart(); err != nil {
						message = fmt.Sprintf("error restarting provider: %v", err)
					} else {
						continue
					}
				} else {
					message = fmt.Sprintf("process exited with code %d", exitErr.ExitCode())
				}
			} else {
				message = fmt.Sprintf("failed to complete process: %v", err)
			}
		} else {
			status = PsCompleted
			message = "process completed"
		}

		break
	}

	p.finalize(status, message)

	return p.Info.Status
}

func (p *Provider) Init() error {
	p.init.Lock()
	defer p.init.Unlock()

	if p.terminated.Load() {
		return errors.New("provider is terminated")
	}

	if err := p.Process.Init(); err != nil {
		return err
	}

	p.Info.Status = PsRunning

	go func() {
		if err := p.Process.Read(); err != nil {
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

	log.Printf("process %d finalized with status %s: %s\n",
		p.Info.Pid, p.Info.Status, p.Info.StatusMessage)

	p.Process.Reset("provider terminated")

	return status, true
}

func (p *Provider) restart() error {
	pr, err := process.Start(p.Cmd)
	if err != nil {
		return err
	}

	p.Process = pr

	if err := p.Init(); err != nil {
		return err
	}

	return nil
}
