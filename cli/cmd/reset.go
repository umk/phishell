package cmd

import (
	"context"
	"fmt"

	"github.com/umk/phishell/thread"
	"github.com/umk/phishell/util/execx"
)

type ResetCommand struct {
	context *Context
}

func (c *ResetCommand) Execute(ctx context.Context, args execx.Arguments) error {
	c.context.session.History = thread.Reset(c.context.session.History)

	fmt.Println("OK")

	return nil
}

func (k *ResetCommand) Info() string {
	return "reset: reset chat history"
}
