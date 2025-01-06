package cmd

import (
	"time"

	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/cli/tool/host"
	"github.com/umk/phishell/util/execx"
)

type Context struct {
	session *session.Session

	commands map[string]Command

	jobs map[int]*backgroundJob
}

type backgroundJob struct {
	args execx.Arguments

	process *host.ToolProcess // TODO: make a weak pointer
	info    *host.ToolProcessInfo

	startedAt time.Time
}

func NewContext(session *session.Session, debug bool) *Context {
	context := &Context{
		session: session,

		commands: make(map[string]Command),
		jobs:     make(map[int]*backgroundJob),
	}

	context.commands["attach"] = &AttachCommand{context: context}
	context.commands["cd"] = &CdCommand{context: context}
	context.commands["export"] = &ExportCommand{}
	context.commands["help"] = &HelpCommand{context: context}
	context.commands["inbox"] = &InboxCommand{context: context}
	context.commands["jobs"] = &JobsCommand{context: context}
	context.commands["kill"] = &KillCommand{context: context}
	context.commands["push"] = &PushCommand{context: context}
	context.commands["pwd"] = &PwdCommand{}
	context.commands["reset"] = &ResetCommand{context: context}

	if debug {
		context.commands["history"] = &HistoryCommand{context: context}
	}

	return context
}

func (c *Context) Command(name string) (Command, bool) {
	cmd, ok := c.commands[name]

	return cmd, ok
}

func (c *Context) refreshJobs() {
	current := make(map[int]*backgroundJob)
	var failedOrCompl struct {
		pid int
		job *backgroundJob
	}

	for pid, job := range c.jobs {
		if job.info.Status == host.TsCompleted || job.info.Status == host.TsFailed {
			// Leave only the most recent failed or completed job for diagnostic purposes
			if failedOrCompl.job == nil || failedOrCompl.job.startedAt.Compare(job.startedAt) < 0 {
				failedOrCompl.pid, failedOrCompl.job = pid, job
			}
		} else {
			current[pid] = job
		}
	}

	if failedOrCompl.job != nil {
		current[failedOrCompl.pid] = failedOrCompl.job
	}

	c.jobs = current
}
