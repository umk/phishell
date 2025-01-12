package cmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os/exec"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/msg"
	"github.com/umk/phishell/response"
	"github.com/umk/phishell/util/execx"
)

type PushCommand struct {
	context *Context
}

func (c *PushCommand) Execute(ctx context.Context, args execx.Arguments) error {
	if len(args) == 0 {
		return c.pushPrevious(ctx)
	}

	return c.pushExec(ctx, args)
}

func (c *PushCommand) Info() string {
	return "push <cmd>: run non-interactive command and push result to chat history"
}

func (c *PushCommand) pushExec(ctx context.Context, args execx.Arguments) error {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	logger := execx.Log(cmd, 15000)

	exitCode, err := execx.Run(cmd)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("process exited with code %d", exitCode)
	}

	processOut, err := logger.Output()
	if err != nil {
		return err
	}

	cr := bootstrap.GetClient(ctx)
	outputStr, err := response.GetExecOutput(ctx, cr, &response.ExecOutputParams{
		CommandLine: args.String(),
		ExitCode:    exitCode,
		Output:      processOut,
	})
	if err != nil {
		return err
	}

	message, err := msg.FormatPushMessage(&msg.PushMessageParams{
		CommandLine: args.String(),
		Output:      outputStr,
	})
	if err != nil {
		return err
	}

	c.context.session.History.Pending = append(c.context.session.History.Pending, message)

	return nil
}

func (c *PushCommand) pushPrevious(ctx context.Context) error {
	session := c.context.session

	if session.PreviousOut == nil {
		return errors.New("no previous output to push")
	}

	cr := bootstrap.GetClient(ctx)
	outputStr, err := response.GetExecOutput(ctx, cr, &response.ExecOutputParams{
		CommandLine: session.PreviousOut.CommandLine,
		ExitCode:    session.PreviousOut.ExitCode,
		Output:      session.PreviousOut.Output,
	})
	if err != nil {
		return err
	}

	message, err := msg.FormatPushMessage(&msg.PushMessageParams{
		CommandLine: session.PreviousOut.CommandLine,
		Output:      outputStr,
	})
	if err != nil {
		return err
	}

	session.History.Pending = append(session.History.Pending, message)
	session.PreviousOut = nil

	fmt.Println("OK")

	return nil
}
