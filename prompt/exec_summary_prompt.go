package prompt

import (
	"context"
	_ "embed"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/client"
	"github.com/umk/phishell/prompt/msg"
	"github.com/umk/phishell/util/execx"
)

type ExecSummaryPromptParams struct {
	Client      *client.Ref
	CommandLine string
	ExitCode    int
	Output      execx.ProcessOutput
}

func PromptExecSummary(ctx context.Context, params *ExecSummaryPromptParams) (*Completion, error) {
	output, tail, err := params.Output.Get()
	if err != nil {
		return nil, err
	}

	message, err := msg.FormatExecSummaryReqMessage(
		&msg.ExecSummaryReqMessageParams{
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

	c, err := params.Client.Completion(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(params.Client.Model(client.Tier1)),
		TopP:     openai.F(0.25),
	})
	if err != nil {
		return nil, err
	}

	return (*Completion)(c), nil
}
