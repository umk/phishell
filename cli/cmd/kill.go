package cmd

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/util/execx"
)

type KillCommand struct {
	context *Context
}

func (c *KillCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) != 1 {
		return ErrInvalidArgs
	}
	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return ErrInvalidArgs
	}

	n := slices.IndexFunc(c.context.providers, func(provider *providerRef) bool {
		return provider.info.Pid == pid
	})
	if n == -1 {
		return fmt.Errorf("no such provider with PID %d", pid)
	}

	p := c.context.providers[n]
	if p.process != nil {
		p.process.Terminate(provider.PsCompleted, "terminated by user")
		p.process = nil
	}

	c.context.providers.refresh()

	fmt.Println("OK")

	return nil
}

func (c *KillCommand) Usage() []string {
	return []string{"kill [pid]"}
}

func (k *KillCommand) Info() []string {
	return []string{"kill a tools provider with the given PID"}
}
