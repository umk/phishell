package prompt

import (
	"context"
	_ "embed"

	"github.com/openai/openai-go"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt/client"
	"github.com/umk/phishell/util/execx"
)

type ExecSummaryPromptParams struct {
	CommandLine string
	ExitCode    int
	Output      execx.ProcessOutput
}

func PromptExecSummary(ctx context.Context, params *ExecSummaryPromptParams) (string, error) {
	output, tail, err := params.Output.Get()
	if err != nil {
		return "", err
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
		return "", err
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(message),
	}

	app := bootstrap.GetApp(ctx)

	cl := client.Get(app.PrimaryClient())

	c, err := cl.Completion(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(cl.GetModel(ctx, client.Tier1)),
		TopP:     openai.F(0.25),
	})
	if err != nil {
		return "", err
	}

	content, err := getCompletionContent(c)
	if err != nil {
		return "", err
	}

	return content, nil
}
