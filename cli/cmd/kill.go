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
		return fmt.Errorf("usage: kill [pid]")
	}
	pidStr := args[0]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return err
	}

	n := slices.IndexFunc(c.context.jobs, func(job *backgroundJob) bool {
		return job.info.Pid == pid
	})
	if n == -1 {
		return fmt.Errorf("no such job: %d", pid)
	}

	job := c.context.jobs[n]
	if job.process != nil {
		job.process.Terminate(host.PsCompleted, "terminated by user")
		job.process = nil
	}

	c.context.refreshJobs()

	fmt.Println("OK")

	return nil
}

func (k *KillCommand) Info() string {
	return "kill [pid]: kill a background process with the given PID"
}
