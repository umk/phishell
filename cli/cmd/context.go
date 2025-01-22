package cmd

import (
	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/tool/host"
	"github.com/umk/phishell/util/execx"
)

type Context struct {
	session *session.Session

	commands map[string]Command

	jobs []*backgroundJob
}

type backgroundJob struct {
	args execx.Arguments

	process *host.Provider
	info    *host.ProviderInfo
}

func NewContext(session *session.Session, debug bool) *Context {
	context := &Context{
		session: session,

		commands: make(map[string]Command),
	}

	context.commands["attach"] = &AttachCommand{context: context}
	context.commands["cd"] = &CdCommand{context: context}
	context.commands["export"] = &ExportCommand{}
	context.commands["help"] = &HelpCommand{context: context}
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
	current := make([]*backgroundJob, 0, len(c.jobs))
	for _, bj := range c.jobs {
		if bj.info.Status == host.PsRunning {
			current = append(current, bj)
		}
	}

	c.jobs = current
}
