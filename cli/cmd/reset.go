package cmd

import (
	"context"
	"fmt"

	"github.com/umk/phishell/util/execx"
)

type ResetCommand struct {
	context *Context
}

func (c *ResetCommand) Execute(ctx context.Context, args execx.Arguments) error {
	c.context.session.History = c.context.session.History.Reset()

	fmt.Println("OK")

	return nil
}

func (k *ResetCommand) Info() string {
	return "reset: reset chat history"
}
