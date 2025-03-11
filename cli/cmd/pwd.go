package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/umk/phishell/util/execx"
)

type PwdCommand struct{}

func (c *PwdCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) > 0 {
		return ErrInvalidArgs
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(workingDir)
	return nil
}

func (c *PwdCommand) Usage() []string {
	return []string{"pwd"}
}

func (p *PwdCommand) Info() []string {
	return []string{"print the current directory"}
}
