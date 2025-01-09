package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/cmd"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/cli/thread"
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

func NewCli(debug bool) *Cli {
	mode := PrCommand

	cli := &Cli{mode: mode}

	cli.session = session.NewSession()
	cli.commands = cmd.NewContext(cli.session, debug)

	return cli
}

func (c *Cli) Init(ctx context.Context) error {
	app := bootstrap.GetApp(ctx)

	if err := os.Chdir(app.Config.Dir); err != nil {
		return err
	}

	cr := bootstrap.GetPrimaryClient(ctx)
	message, err := msg.FormatSystemMessage(&msg.SystemMessageParams{
		Prompt: cr.Config.Prompt,
		OS:     runtime.GOOS,
	})
	if err != nil {
		return err
	}

	c.session.History = &thread.History{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(message),
		},
	}

	return nil
}

func (c *Cli) Run(ctx context.Context) error {
	defer c.session.Host.Close()

	cancelThisContext := cancelOnSigTerm()

	for {
		ctx := cancelThisContext(ctx)

		if err := c.processPrompt(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				processCliError(err)
			}
		}
	}

	return nil
}

func (c *Cli) getClient(ctx context.Context) *bootstrap.ClientRef {
	n := int(c.mode) - int(PrChat)
	if n < 0 {
		panic("prompt is not in a chat mode")
	}

	app := bootstrap.GetApp(ctx)

	return app.Clients[n]
}

func processCliError(err error) {
	if errorsx.IsCanceled(err) {
		// Do nothing
	} else {
		termx.Error.Println(err)
	}
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
