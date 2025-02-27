package cli

import (
	"context"
	"strings"

	"github.com/umk/phishell/util/termx"
)

func (c *Cli) processPrompt(ctx context.Context) error {
	line, err := termx.ReadPrompt(ctx, &promptChat{cli: c})
	if err != nil {
		return err
	}

	// Handle empty content
	content := strings.TrimSpace(line)
	if content == "" {
		return nil
	}

	switch c.mode {
	case PrCommand:
		return c.processCommand(ctx, content)
	case PrChat:
		client := c.getClient(ctx)
		return c.session.Post(ctx, client, content)
	}

	return nil
}
