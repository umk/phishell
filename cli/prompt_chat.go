package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/umk/phishell/bootstrap"
)

type promptChat struct {
	cli *Cli
}

func (p *promptChat) getPrompt(ctx context.Context, mode PromptMode) string {
	switch mode {
	case PrCommand:
		return p.getCommandPrompt()
	default:
		return p.getChatPrompt(ctx)
	}
}

func (p *promptChat) getHint(ctx context.Context, mode PromptMode) string {
	if mode == PrCommand {
		return ""
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	homeDir, err := os.UserHomeDir()

	if err == nil && dir == homeDir {
		return "~"
	} else if dir == "/" {
		return "/"
	} else {
		dirName := filepath.Base(dir)
		return dirName + "/"
	}
}

func (p *promptChat) getNextMode(ctx context.Context, current PromptMode) PromptMode {
	app := bootstrap.GetApp(ctx)

	max := int(PrChat) + len(app.Clients)

	return PromptMode((int(current) + 1) % max)
}

func (p *promptChat) getChatPrompt(ctx context.Context) string {
	client := p.cli.getClient(ctx)

	return fmt.Sprintf("%s >>> ", client.Config.Profile)
}

func (p *promptChat) getCommandPrompt() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	var base string

	homeDir, err := os.UserHomeDir()

	if err == nil && dir == homeDir {
		base = "~"
	} else if dir == "/" {
		base = "/"
	} else {
		dirName := filepath.Base(dir)
		base = dirName
	}

	return base + " $ "
}
