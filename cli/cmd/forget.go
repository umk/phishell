package cmd

import (
	"context"

	"github.com/umk/phishell/util/execx"
)

type ForgetCommand struct {
	context *Context
}

func (c *ForgetCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) != 1 {
		return ErrInvalidArgs
	}

	return nil
}

func (c *ForgetCommand) Usage() []string {
	return []string{"forget [batch]"}
}

func (c *ForgetCommand) Info() []string {
	return []string{"forget the previously learned batch"}
}
