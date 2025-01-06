package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/umk/phishell/util/execx"
)

type PwdCommand struct{}

func (c *PwdCommand) Execute(ctx context.Context, args execx.Arguments) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(workingDir)
	return nil
}

func (p *PwdCommand) Info() string {
	return "pwd: print the current directory"
}
