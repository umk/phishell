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

func (p *promptChat) GetPrompt(ctx context.Context) string {
	switch p.cli.mode {
	case PrCommand:
		return p.getCommandPrompt()
	default:
		return p.getChatPrompt(ctx)
	}
}

func (p *promptChat) GetHint(ctx context.Context) string {
	if p.cli.mode == PrCommand {
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

func (p *promptChat) CycleMode(ctx context.Context) {
	clients := bootstrap.GetClients(ctx)

	max := int(PrChat) + len(clients)

	p.cli.mode = PromptMode((int(p.cli.mode) + 1) % max)
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
