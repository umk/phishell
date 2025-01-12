package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/umk/phishell/tool/host"
	"github.com/umk/phishell/util/execx"
)

type KillCommand struct {
	context *Context
}

func (c *KillCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: kill [pid]")
	}
	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return err
	}

	cmd, ok := c.context.jobs[pid]
	if !ok {
		return fmt.Errorf("no such job: %d", pid)
	}

	cmd.process.Terminate(ctx, host.TsCompleted, "terminated by user")

	c.context.refreshJobs()

	fmt.Println("OK")

	return nil
}

func (k *KillCommand) Info() string {
	return "kill [pid]: kill a background process with the given PID"
}
