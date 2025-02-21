package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/umk/phishell/util/execx"
)

type AttachCommand struct {
	context *Context
}

func (c *AttachCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: %s", c.Usage())
	}

	cmd := execx.Cmd{
		Cmd:  args[0],
		Args: args[1:],
		Env:  append(os.Environ(), "PHI_SHELL=1"),
	}

	p, err := c.context.session.Host.Execute(&cmd)
	if err != nil {
		return err
	}

	pr := &providerRef{
		args:    args,
		process: p,
		info:    p.Info(),
	}

	c.context.providers = append(c.context.providers, pr)

	fmt.Printf("started provider [%d]\n", pr.info.Pid)

	return nil
}

func (c *AttachCommand) Usage() string {
	return "attach [cmd]"
}

func (p *AttachCommand) Info() string {
	return "run tools provider in background"
}
