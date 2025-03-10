package cmd

import (
	"context"
	"errors"
	"strings"

	"github.com/umk/phishell/util/execx"
)

type Command interface {
	Execute(ctx context.Context, args execx.Arguments) error

	Usage() []string
	Info() []string
}

func getUsageError(c Command) error {
	var sb strings.Builder

	sb.WriteString("usage: ")
	for i, usage := range c.Usage() {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(usage)
	}

	return errors.New(sb.String())
}
