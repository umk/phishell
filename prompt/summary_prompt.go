package prompt

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/prompt/msg"
)

type SummaryPromptParams struct {
	Messages []openai.ChatCompletionMessageParamUnion
}

func PromptSummary(ctx context.Context, params *SummaryPromptParams) (*Completion, error) {
	m, err := msg.FormatSummaryReqMessage(&msg.SummaryReqMessageParams{})
	if err != nil {
		return nil, err
	}

	n := len(params.Messages)

	messages := make([]openai.ChatCompletionMessageParamUnion, n, n+1)

	copy(messages, params.Messages)
	messages = append(messages, openai.UserMessage(m))

	c, err := client.Default.Completion(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    client.Default.Model(client.Tier1),
		TopP:     param.NewOpt(0.25),
	})
	if err != nil {
		return nil, err
	}

	return (*Completion)(c), nil
}
