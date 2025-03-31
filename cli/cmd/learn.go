package cmd

import (
	"context"

	"github.com/umk/phishell/util/execx"
)

type LearnCommand struct {
	context *Context
}

func (c *LearnCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) < 1 {
		return ErrInvalidArgs
	}

	return nil
}

func (c *LearnCommand) Usage() []string {
	return []string{"learn [pattern] ..."}
}

func (c *LearnCommand) Info() []string {
	return []string{"find and learn documents by pattern"}
}
