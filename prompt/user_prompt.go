package prompt

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/client"
)

type UserPromptParams struct {
	Client   *client.Ref
	Messages []openai.ChatCompletionMessageParamUnion
	Tools    []openai.ChatCompletionToolParam
}

func PromptUser(ctx context.Context, params *UserPromptParams) (*Completion, error) {
	c, err := params.Client.Completion(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(params.Messages),
		Model:    openai.F(params.Client.Model(client.Tier1)),
		Tools:    openai.F(params.Tools),
		TopP:     openai.F(0.25),
	})
	if err != nil {
		return nil, err
	}

	return (*Completion)(c), nil
}
