package prompt

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/prompt/client"
)

type UserPromptParams struct {
	Client   *bootstrap.ClientRef
	Messages []openai.ChatCompletionMessageParamUnion
	Tools    []openai.ChatCompletionToolParam
}

func PromptUser(ctx context.Context, params *UserPromptParams) (*openai.ChatCompletion, error) {
	cl := client.Get(params.Client)

	return cl.Completion(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(params.Messages),
		Model:    openai.F(cl.GetModel(ctx, client.Tier1)),
		Tools:    openai.F(params.Tools),
		TopP:     openai.F(0.25),
	})
}
