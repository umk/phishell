package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/umk/phishell/cli/cmd"
	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/config"
	"github.com/umk/phishell/util/errorsx"
	"github.com/umk/phishell/util/termx"
)

type PromptMode int

const (
	PrCommand PromptMode = iota
	PrChat
)

type Cli struct {
	mode PromptMode

	session  *session.Session
	commands *cmd.Context
}

func NewCli() *Cli {
	cli := &Cli{mode: PrCommand}

	cli.session = session.NewSession()
	cli.commands = cmd.NewContext(cli.session)

	return cli
}

func (c *Cli) Init(ctx context.Context) error {
	if err := os.Chdir(config.Config.Dir); err != nil {
		return err
	}

	return nil
}

func (c *Cli) Run(ctx context.Context) error {
	defer c.session.Host.Close()

	cancelThisContext := cancelOnSigTerm()

	for {
		termx.MD.Init()

		ctx := cancelThisContext(ctx)

		if err := c.processPrompt(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			} else if errorsx.IsCanceled(err) {
				// Do nothing
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func (c *Cli) getClient() *client.Ref {
	n := int(c.mode) - int(PrChat)
	if n < 0 {
		panic("prompt is not in a chat mode")
	}

	return client.ChatProfiles[n]
}

func cancelOnSigTerm() func(context.Context) context.Context {
	// Set up signal handling for Ctrl+C
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGTERM)

	var cancelCur func()

	go func() {
		for range termChan {
			cur := cancelCur
			if cur != nil {
				cur()
			}
		}
	}()

	return func(ctx context.Context) context.Context {
		ctx, cancel := context.WithCancel(ctx)
		cancelCur = cancel

		return ctx
	}
}
