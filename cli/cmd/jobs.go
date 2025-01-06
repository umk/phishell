package cmd

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/umk/phishell/util/execx"
)

type JobsCommand struct {
	context *Context
}

func (c *JobsCommand) Execute(ctx context.Context, args execx.Arguments) error {
	c.context.refreshJobs()

	if len(c.context.jobs) == 0 {
		fmt.Println("no running or stopped jobs")
		return nil
	}

	var pids []int
	for pid := range c.context.jobs {
		pids = append(pids, pid)
	}

	slices.Sort(pids)

	for _, pid := range pids {
		job := c.context.jobs[pid]

		var status strings.Builder

		status.WriteString(job.info.Status.String())
		if job.info.StatusMessage != "" {
			status.WriteString("; ")
			status.WriteString(job.info.StatusMessage)
		}

		fmt.Printf("[%d] %s (%s)\n", pid, job.args, status.String())
	}
	return nil
}

func (j *JobsCommand) Info() string {
	return "jobs: list background processes"
}
