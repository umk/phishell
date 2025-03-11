package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/umk/phishell/cli/cmd"
	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/config"
	"github.com/umk/phishell/util/execx"
	"github.com/umk/phishell/util/termx"
)

// stdoutWrapper wraps an io.Writer to track if the last character was a newline
type stdoutWrapper struct {
	w         io.Writer
	lastWasNL bool
}

func (s *stdoutWrapper) Write(p []byte) (n int, err error) {
	n, err = s.w.Write(p)
	if len(p) > 0 {
		s.lastWasNL = p[len(p)-1] == '\n' || p[len(p)-1] == '\r'
	}
	return n, err
}

func (c *Cli) processCommand(ctx context.Context, content string) error {
	piped, err := execx.Parse(content)
	if err != nil {
		return err
	}

	for _, args := range piped {
		if len(args) == 0 {
			if len(piped) > 1 {
				return errors.New("empty command")
			}

			return nil
		}

		execCmd, err := args.Cmd()
		if err != nil {
			return err
		}

		if b, ok := c.commands.Command(execCmd.Cmd); ok {
			if len(piped) > 1 {
				return errors.New("cannot pipe built-in command")
			}

			if len(execCmd.Env) > 0 {
				return errors.New("cannot assign environment variables when calling built-in command")
			}

			if err := b.Execute(ctx, execCmd.Args); err != nil {
				if err == cmd.ErrInvalidArgs {
					printUsageError(b)
					return nil
				} else {
					return err
				}
			}

			return nil
		}
	}

	return c.processExternalCommand(ctx, piped)
}

func (c *Cli) processExternalCommand(ctx context.Context, piped []execx.Arguments) error {
	if len(piped) == 0 {
		return nil
	}

	// Reset the previous output
	c.session.PreviousOut = nil

	cmds := make([]*exec.Cmd, len(piped))
	for i, p := range piped {
		cmd, err := p.Cmd()
		if err != nil {
			return err
		}

		cmds[i] = cmd.CommandContext(ctx)
	}

	wrapper := &stdoutWrapper{w: os.Stdout, lastWasNL: true}
	if err := execx.Pipe(cmds, os.Stdin, wrapper, wrapper); err != nil {
		return err
	}

	logger := execx.Log(cmds[len(cmds)-1], config.Config.OutputBufSize)

	exitCode, err := execx.RunPipe(cmds)

	// Print newline if output didn't end with newline
	if !wrapper.lastWasNL {
		fmt.Println()
	}

	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("exit code %d", exitCode)
	}

	// Set the previous output
	output, err := logger.Output()
	if err == nil {
		c.session.PreviousOut = &session.PreviousOut{
			CommandLine: piped[len(piped)-1].String(),
			ExitCode:    0,
			Output:      output,
		}
	} else if config.Config.Debug {
		termx.Muted.Printf("(error) %v\n", err)
	}

	return nil
}

func printUsageError(c cmd.Command) {
	const UsagePrefix = "usage: "
	padding := strings.Repeat(" ", utf8.RuneCountInString(UsagePrefix))

	for i, usage := range c.Usage() {
		if i == 0 {
			fmt.Print(UsagePrefix)
		} else {
			fmt.Print(padding)
		}
		fmt.Println(usage)
	}
}
