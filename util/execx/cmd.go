package execx

import (
	"context"
	"os"
	"os/exec"
)

type Cmd struct {
	Env  []string
	Cmd  string
	Args []string
}

func (c *Cmd) Command() *exec.Cmd {
	cmd := exec.Command(c.Cmd, c.Args...)
	cmd.Env = c.Env

	return cmd
}

func (c *Cmd) CommandContext(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.Cmd, c.Args...)
	cmd.Env = append(os.Environ(), cmd.Env...)

	return cmd
}
