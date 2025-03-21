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
	if len(args) > 0 {
		return ErrInvalidArgs
	}

	c.context.session.History = new(thread.History)

	fmt.Println("OK")

	return nil
}

func (c *ResetCommand) Usage() []string {
	return []string{"reset"}
}

func (c *ResetCommand) Info() []string {
	return []string{"reset chat history"}
}
