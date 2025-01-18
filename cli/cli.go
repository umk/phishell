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
	"github.com/umk/phishell/cli/session"
	"github.com/umk/phishell/msg"
	"github.com/umk/phishell/thread"
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
	config := bootstrap.GetConfig(ctx)

	if err := os.Chdir(config.Dir); err != nil {
		return err
	}

	var script string
	if bootstrap.IsScript(ctx) {
		s, err := readScript(config.Startup.Script)
		if err != nil {
			return err
		}

		script = s
	}

	cr := bootstrap.GetClient(ctx)
	message, err := msg.FormatSystemMessage(&msg.SystemMessageParams{
		Prompt: cr.Config.Prompt,
		Script: script,
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
	defer c.session.Host.Close(ctx)

	cancelThisContext := cancelOnSigTerm()

	if bootstrap.IsScript(ctx) {
		ctx := cancelThisContext(ctx)

		if err := c.processScriptPrompt(ctx); err != nil {
			return err
		}
	}

	for {
		ctx := cancelThisContext(ctx)

		if err := c.processPrompt(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			} else if errorsx.IsCanceled(err) {
				// Do nothing
			} else {
				termx.Error.Println(err)
			}
		}
	}
}

func (c *Cli) getClient(ctx context.Context) *bootstrap.ClientRef {
	n := int(c.mode) - int(PrChat)
	if n < 0 {
		panic("prompt is not in a chat mode")
	}

	clients := bootstrap.GetClients(ctx)

	return clients[n]
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
