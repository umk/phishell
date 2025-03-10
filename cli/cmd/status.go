package cmd

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/umk/phishell/util/execx"
)

type StatusCommand struct {
	context *Context
}

func (c *StatusCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) > 0 {
		return getUsageError(c)
	}

	c.context.providers.refresh()

	if len(c.context.providers) == 0 {
		fmt.Println("no running or stopped providers")
		return nil
	}

	var pids []int
	for pid := range c.context.providers {
		pids = append(pids, pid)
	}

	slices.Sort(pids)

	for _, pid := range pids {
		provider := c.context.providers[pid]

		var status strings.Builder

		status.WriteString(provider.info.Status.String())
		if provider.info.StatusMessage != "" {
			status.WriteString("; ")
			status.WriteString(provider.info.StatusMessage)
		}

		fmt.Printf("[%d] %s (%s)\n", pid, provider.args, status.String())
	}
	return nil
}

func (c *StatusCommand) Usage() []string {
	return []string{"status"}
}

func (c *StatusCommand) Info() []string {
	return []string{"list tool providers"}
}
