package prompt

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/prompt/client"
	"github.com/umk/phishell/prompt/msg"
)

type SummaryPromptParams struct {
	Messages []openai.ChatCompletionMessageParamUnion
}

func PromptSummary(ctx context.Context, params *SummaryPromptParams) (*Completion, error) {
	cl := client.Get(bootstrap.GetDefaultClient(ctx))

	m, err := msg.FormatSummaryReqMessage(&msg.SummaryReqMessageParams{})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return (*Completion)(c), nil
}
