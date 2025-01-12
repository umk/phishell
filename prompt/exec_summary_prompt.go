package prompt

import (
	"context"
	_ "embed"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/msg"
	"github.com/umk/phishell/prompt/client"
	"github.com/umk/phishell/util/execx"
)

type ExecSummaryPromptParams struct {
	Client      *bootstrap.ClientRef
	CommandLine string
	ExitCode    int
	Output      execx.ProcessOutput
}

func PromptExecSummary(ctx context.Context, params *ExecSummaryPromptParams) (*Completion, error) {
	cl := client.Get(params.Client)

	output, tail, err := params.Output.Get()
	if err != nil {
		return nil, err
	}

	message, err := msg.FormatExecSummaryMessage(
		&msg.ExecSummaryMessageParams{
			CommandLine: params.CommandLine,
			ExitCode:    params.ExitCode,
			Output:      output,
			Tail:        tail,
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
