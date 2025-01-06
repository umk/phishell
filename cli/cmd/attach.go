package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/umk/phishell/util/execx"
)

type AttachCommand struct {
	context *Context
}

func (c *AttachCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: attach [cmd]")
	}

	cmd := execx.Cmd{
		Cmd:  args[0],
		Args: args[1:],
		Env:  append(os.Environ(), "PHI_SHELL=1"),
	}

	pr, err := c.context.session.Host.Execute(ctx, &cmd)
	if err != nil {
		return err
	}

	pid := pr.Pid()
	c.context.jobs[pid] = &backgroundJob{
		args:      args,
		process:   pr,
		info:      pr.Info,
		startedAt: time.Now(),
	}

	fmt.Printf("started background job [%d]\n", pid)

	return nil
}

func (p *AttachCommand) Info() string {
	return "attach [cmd]: run background process that provides tools"
}
