package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/umk/phishell/util/execx"
	"github.com/umk/phishell/util/termx"
)

type HelpCommand struct {
	context *Context
}

func (c *HelpCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) > 0 {
		return fmt.Errorf("usage: %s", c.Usage())
	}

	var names []string
	for name := range c.context.commands {
		names = append(names, name)
	}

	slices.Sort(names)

	var commands []string

	for _, name := range names {
		command := c.context.commands[name]
		info := fmt.Sprintf(" - `%s` %s", command.Usage(), command.Info())
		commands = append(commands, info)
	}

	printer := termx.NewPrinter()
	printer.Printf("# Commands\n\n%s\n\nPress the `Tab` key to cycle through the chat profiles and return to the command line.", strings.Join(commands, "\n"))

	return nil
}

func (c *HelpCommand) Usage() string {
	return "help"
}

func (c *HelpCommand) Info() string {
	return "display the help message"
}
