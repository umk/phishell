package cmd

import (
	"context"
	"errors"

	"github.com/umk/phishell/util/execx"
)

var ErrInvalidArgs = errors.New("bad command usage")

type Command interface {
	Execute(ctx context.Context, args execx.Arguments) error

	Usage() []string
	Info() []string
}
