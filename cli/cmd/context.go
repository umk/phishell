package cmd

import (
	"github.com/umk/phishell/cli/session"
)

type Context struct {
	session *session.Session

	commands map[string]Command

	providers providersContext
	documents documentsContext
}

func NewContext(session *session.Session) *Context {
	context := &Context{
		session: session,

		commands: make(map[string]Command),

		documents: makeDocumentsContext(),
	}

	context.commands["attach"] = &AttachCommand{context: context}
	context.commands["cd"] = &CdCommand{context: context}
	context.commands["export"] = &ExportCommand{}
	context.commands["forget"] = &ForgetCommand{context: context}
	context.commands["help"] = &HelpCommand{context: context}
	context.commands["kill"] = &KillCommand{context: context}
	context.commands["learn"] = &LearnCommand{context: context}
	context.commands["push"] = &PushCommand{context: context}
	context.commands["pwd"] = &PwdCommand{}
	context.commands["reset"] = &ResetCommand{context: context}
	context.commands["status"] = &StatusCommand{context: context}

	return context
}

func (c *Context) Command(name string) (Command, bool) {
	cmd, ok := c.commands[name]

	return cmd, ok
}
