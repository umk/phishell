package cmd

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/umk/phishell/tool/host"
	"github.com/umk/phishell/util/execx"
)

type KillCommand struct {
	context *Context
}

func (c *KillCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: %s", c.Usage())
	}
	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return err
	}

	n := slices.IndexFunc(c.context.providers, func(provider *providerRef) bool {
		return provider.info.Pid == pid
	})
	if n == -1 {
		return fmt.Errorf("no such provider: %d", pid)
	}

	provider := c.context.providers[n]
	if provider.process != nil {
		provider.process.Terminate(host.PsCompleted, "terminated by user")
		provider.process = nil
	}

	c.context.providers.refresh()

	fmt.Println("OK")

	return nil
}

func (c *KillCommand) Usage() string {
	return "kill [pid]"
}

func (k *KillCommand) Info() string {
	return "kill a tools provider with the given PID"
}
