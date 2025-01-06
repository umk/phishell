package prompt

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt/client"
)

type SummaryPromptParams struct {
	Messages []openai.ChatCompletionMessageParamUnion
}

func PromptSummary(ctx context.Context, params *SummaryPromptParams) (string, error) {
	app := bootstrap.GetApp(ctx)

	cl := client.Get(app.PrimaryClient())

	m, err := msg.FormatSummaryReqMessage(&msg.SummaryReqMessageParams{})
	if err != nil {
		return "", err
	}

	n := len(params.Messages)

	messages := make([]openai.ChatCompletionMessageParamUnion, n, n+1)

	copy(messages, params.Messages)
	messages = append(messages, openai.UserMessage(m))

	c, err := cl.Completion(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(cl.GetModel(ctx, client.Tier1)),
		TopP:     openai.F(0.25),
	})
	if err != nil {
		return "", err
	}

	summary, err := getCompletionContent(c)
	if err != nil {
		return "", err
	}

	return msg.FormatSummaryMessage(&msg.SummaryMessageParams{
		Summary: summary,
	})
}
