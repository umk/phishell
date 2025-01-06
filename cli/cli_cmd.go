package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/session"
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

		cmd, err := execx.AllocArgs(args)
		if err != nil {
			return err
		}

		if b, ok := c.commands.Command(cmd.Cmd); ok {
			if len(piped) > 1 {
				return errors.New("cannot pipe built-in command")
			}

			if len(cmd.Env) > 0 {
				return errors.New("cannot assign environment variables when calling built-in command")
			}

			return b.Execute(ctx, cmd.Args)
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
		cmd, err := execx.AllocArgs(p)
		if err != nil {
			return err
		}

		cmds[i] = cmd.CommandContext(ctx)
	}

	wrapper := &stdoutWrapper{w: os.Stdout, lastWasNL: true}
	if err := execx.Pipe(cmds, os.Stdin, wrapper, wrapper); err != nil {
		return err
	}

	logger := execx.Log(cmds[len(cmds)-1], 15000)

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
	} else if bootstrap.IsDebug(ctx) {
		termx.Muted.Printf("(error) %s\n", err)
	}

	return nil
}
