package process

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/provider"
	"github.com/umk/phishell/util/execx"
)

type Process struct {
	mu sync.Mutex

	cmd  *exec.Cmd
	wait struct {
		once sync.Once // synchronizes calls to Cmd.Wait
		err  error     // error returned from Cmd.Wait
	}

	stdin  io.Writer
	stdout io.Reader

	// A scanner over the process' messages read from Stdout.
	scanner *bufio.Scanner

	// Mapping from process' tools call ID to a channel, that accepts either
	// the response or error.
	requests map[string]chan<- any

	tools map[string]openai.ChatCompletionToolParam
}

func Start(c *execx.Cmd) (*Process, error) {
	cmd := c.Command()

	pr, err := New(cmd)
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	return pr, nil
}

func New(cmd *exec.Cmd) (*Process, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe Stdin: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to pipe Stdout: %w", err)
	}

	return &Process{
		cmd: cmd,

		stdin:  stdin,
		stdout: stdout,

		requests: make(map[string]chan<- any),

		tools: make(map[string]openai.ChatCompletionToolParam),
	}, nil
}

func (p *Process) Tools() map[string]openai.ChatCompletionToolParam {
	return p.tools
}

func (p *Process) Cmd() *exec.Cmd {
	return p.cmd
}

// WaitOnce waits for the command completion. Safe for calling multiple times.
func (p *Process) WaitOnce() error {
	p.wait.once.Do(func() {
		p.wait.err = p.cmd.Wait()
	})

	return p.wait.err
}

// Post sends the request to a tools provider and waits for a response.
func (p *Process) Post(req *provider.ToolRequest) (*provider.ToolResponse, error) {
	ch, err := p.requestSend(req)
	if err != nil {
		return nil, err
	}

	defer p.requestClose(req.CallID)

	timeout := 15 * time.Second

	select {
	case res, ok := <-ch:
		if !ok {
			return nil, errors.New("request canceled")
		}
		switch v := res.(type) {
		case *provider.ToolResponse:
			return v, nil
		case error:
			return nil, v
		default:
			panic("bad message type")
		}
	case <-time.After(timeout):
		return nil, fmt.Errorf("handler timed out after %s", timeout)
	}
}

func (p *Process) Reset(reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ch := range p.requests {
		ch <- fmt.Errorf("request canceled: %s", reason)
	}

	p.requests = make(map[string]chan<- any)
}
