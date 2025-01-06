package thread

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/prompt"
	"github.com/umk/phishell/util/termx"
)

func (t *Thread) processChatMessage(ctx context.Context, response openai.ChatCompletionMessage, sizeToks int64) error {
	if response.Refusal != "" {
		return fmt.Errorf("refused: %s", response.Refusal)
	}

	t.history.Messages = append(t.history.Messages, response)

	if sizeToks > int64(t.client.Config.CompactionToks) {
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
	}

	fmt.Println(response.Content)

	return nil
}
