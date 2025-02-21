package cmd

import (
	"context"

	"github.com/umk/phishell/util/execx"
)

type Command interface {
	Execute(ctx context.Context, args execx.Arguments) error
	Usage() string
	Info() string
}
