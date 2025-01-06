package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/umk/phishell/util/execx"
)

type HelpCommand struct {
	context *Context
}

func (c *HelpCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) > 0 {
		return fmt.Errorf("usage: help")
	}

	var names []string

	for name := range c.context.commands {
		names = append(names, name)
	}

	slices.Sort(names)

	var commands []string

	for _, name := range names {
		command := c.context.commands[name]

		info := fmt.Sprintf(" - %s", command.Info())

		commands = append(commands, info)
	}

	fmt.Println(strings.Join(commands, "\n"))

	return nil
}

func (c *HelpCommand) Info() string {
	return "help: display the help message"
}
