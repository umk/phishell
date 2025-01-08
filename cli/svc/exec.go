package svc

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt"
	"github.com/umk/phishell/util/execx"
)

type ExecOutputParams struct {
	CommandLine string
	ExitCode    int
	Output      execx.ProcessOutput
}

func GetExecOutput(ctx context.Context, params *ExecOutputParams) (string, error) {
	outputStr, tail, err := params.Output.Get()
	if err != nil {
		return "", fmt.Errorf("invalid output: %w", err)
	}

	if params.ExitCode == 0 {
		if s, ok := getJSONObjectOrArray(outputStr); ok {
			outputStr = s
		}
	}

	if utf8.RuneCountInString(outputStr) < 2500 {
		return msg.FormatExecResponseMessage(&msg.ExecResponseMessageParams{
			ExitCode: params.ExitCode,
			Output:   outputStr,
			Summary:  false,
		})
	}

	summaryCompl, err := prompt.PromptExecSummary(ctx, &prompt.ExecSummaryPromptParams{
		CommandLine: params.CommandLine,
		ExitCode:    params.ExitCode,
		Output:      params.Output,
	})
	if err != nil {
		return "", err
	}

	summary, err := summaryCompl.Content()
	if err != nil {
		return "", err
	}

	return msg.FormatExecResponseMessage(&msg.ExecResponseMessageParams{
		ExitCode: params.ExitCode,
		Output:   summary,
		Summary:  tail,
	})
}
