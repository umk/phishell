package cli

import "context"

type promptScript struct {
	cli *Cli
}

func (p *promptScript) getPrompt(ctx context.Context, mode PromptMode) string {
	return ">>> "
}

func (p *promptScript) getHint(ctx context.Context, mode PromptMode) string {
	return "provide instructions or press Enter to continue"
}

func (p *promptScript) getNextMode(ctx context.Context, current PromptMode) PromptMode {
	return p.cli.mode // preserve current mode
}
