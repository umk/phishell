package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/umk/phishell/util/execx"
	"github.com/umk/phishell/util/fsx"
)

type AttachCommand struct {
	context *Context
}

func (c *AttachCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) < 1 {
		return ErrInvalidArgs
	}

	var cmd execx.Cmd
	var err error

	if args[0] == "js" {
		cmd, err = createAttachJS(args)
		if err != nil {
			return err
		}
	} else {
		cmd = createAttachGeneric(args)
	}

	p, err := c.context.session.Host.Execute(&cmd)
	if err != nil {
		return err
	}

	pr := &providerRef{
		args:    args,
		process: p,
		info:    p.Info,
	}

	c.context.providers = append(c.context.providers, pr)

	fmt.Printf("started provider [%d]\n", pr.info.Pid)

	return nil
}

func (c *AttachCommand) Usage() []string {
	return []string{
		"attach [cmd]",
		"attach js [path]",
	}
}

func (p *AttachCommand) Info() []string {
	return []string{
		"run tools provider in background",
		"attach Node.js package that implements tools",
	}
}

func createAttachJS(args execx.Arguments) (execx.Cmd, error) {
	wd, err := os.Getwd()
	if err != nil {
		return execx.Cmd{}, fmt.Errorf("unable to get current working directory: %w", err)
	}

	var dir string

	switch len(args) {
	case 1:
		dir = wd
	case 2:
		dir = fsx.Resolve(wd, args[1])
	default:
		return execx.Cmd{}, ErrInvalidArgs
	}

	executable, err := os.Executable()
	if err != nil {
		return execx.Cmd{}, fmt.Errorf("unable to determine executable location: %w", err)
	}

	jsPath := filepath.Join(filepath.Dir(executable), "phishell-js.mjs")
	return execx.Cmd{
		Cmd:  "node",
		Args: []string{jsPath, "--", "serve"},
		Dir:  dir,
	}, nil
}

func createAttachGeneric(args execx.Arguments) execx.Cmd {
	return execx.Cmd{
		Cmd:  args[0],
		Args: args[1:],
		Env:  append(os.Environ(), "PHI_SHELL=1"),
	}
}
