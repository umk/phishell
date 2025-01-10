package prompt

import (
	"context"
	_ "embed"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt/client"
)

type HttpSummaryPromptParams struct {
	Client  *bootstrap.ClientRef
	Url     string
	Status  string
	Headers map[string][]string
	Body    string
}

func PromptHttpSummary(ctx context.Context, params *HttpSummaryPromptParams) (*Completion, error) {
	cl := client.Get(params.Client)

	message, err := msg.FormatHttpSummaryMessage(
		&msg.HttpSummaryMessageParams{
			Url:     params.Url,
			Status:  params.Status,
			Headers: params.Headers,
			Body:    params.Body,
		},
	)
	if err != nil {
		return nil, err
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(message),
	}

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
