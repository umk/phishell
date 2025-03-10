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
	Dir  string
}

func (c *Cmd) Command() *exec.Cmd {
	cmd := exec.Command(c.Cmd, c.Args...)
	cmd.Env = c.Env
	cmd.Dir = c.Dir

	return cmd
}

func (c *Cmd) CommandContext(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.Cmd, c.Args...)
	cmd.Env = append(os.Environ(), cmd.Env...)

	return cmd
}
