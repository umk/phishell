package thread

import (
	"context"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/prompt"
	"github.com/umk/phishell/util/termx"
)

func (t *Thread) compactHistory(ctx context.Context) error {
	if t.history.Toks < int64(t.client.Config.CompactionToks) {
		return nil
	}

	if bootstrap.IsDebug(ctx) {
		termx.Muted.Println("(compaction)")
	}

	summary, err := prompt.PromptSummary(ctx, &prompt.SummaryPromptParams{
		Messages: t.history.Messages,
	})
	if err != nil {
		return err
	}

	if bootstrap.IsDebug(ctx) {
		termx.Muted.Println(summary)
	}

	history := t.history.Reset()
	history.Pending = append(history.Pending, summary)

	t.history = history

	return nil
}
