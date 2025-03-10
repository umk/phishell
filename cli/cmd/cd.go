package cmd

import (
	"context"
	"os"

	"github.com/umk/phishell/util/execx"
)

type CdCommand struct {
	context *Context
}

func (c *CdCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) > 1 {
		return getUsageError(c)
	}

	var dir string
	if len(args) == 0 {
		dir = os.Getenv("HOME")
	} else {
		dir = args[0]
	}

	dir, err := c.context.session.Resolve(dir)
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	return nil
}

func (c *CdCommand) Usage() []string {
	return []string{"cd <dir>"}
}

func (c *CdCommand) Info() []string {
	return []string{"change the current directory"}
}
