package thread

import (
	"context"

	"github.com/umk/phishell/prompt"
	"github.com/umk/phishell/prompt/client"
	"github.com/umk/phishell/prompt/msg"
	"github.com/umk/phishell/util/termx"
)

func (t *Thread) compactHistory(ctx context.Context) error {
	if t.history.Toks < t.client.Config.CompactionToks {
		return nil
	}

	termx.Muted.Println("Compaction...")

	summaryCompl, err := prompt.PromptSummary(ctx, &prompt.SummaryPromptParams{
		Messages: t.history.Messages,
	})
	if err != nil {
		return err
	}

	summary, err := summaryCompl.Content()
	if err != nil {
		return err
	}

	summaryMsg, err := msg.FormatSummaryMessage(&msg.SummaryMessageParams{
		Summary: summary,
	})
	if err != nil {
		return err
	}

	termx.Muted.Println(summaryMsg)

	c := client.Get(t.client)

	systemToks := int(float32(getSystemMessageSize(t.history)) / c.Samples.BytesPerTok())
	summaryToks := summaryCompl.Toks() // doesn't consider content added by the summary template

	toks := systemToks + summaryToks

	if toks >= t.history.Toks {
		return nil
	}

	history := Reset(t.history)
	history.Pending = append(history.Pending, summaryMsg)
	history.Toks = toks

	t.history = history

	return nil
}
