package svc

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/umk/phishell/cli/msg"
	"github.com/umk/phishell/cli/prompt"
	"github.com/umk/phishell/util/execx"
)

func GetExecOutput(ctx context.Context, commandLine string, exitCode int, output execx.ProcessOutput) (string, error) {
	outputStr, tail, err := output.Get()
	if err != nil {
		return "", fmt.Errorf("invalid output: %w", err)
	}

	if exitCode == 0 {
		if s, ok := getJSONObjectOrArray(outputStr); ok {
			outputStr = s
		}
	}

	if utf8.RuneCountInString(outputStr) < 2500 {
		return msg.FormatExecResponseMessage(&msg.ExecResponseMessageParams{
			ExitCode: exitCode,
			Output:   outputStr,
			Summary:  false,
		})
	}

	summaryCompl, err := prompt.PromptExecSummary(ctx, &prompt.ExecSummaryPromptParams{
		CommandLine: commandLine,
		ExitCode:    exitCode,
		Output:      output,
	})
	if err != nil {
		return "", err
	}

	summary, err := summaryCompl.Content()
	if err != nil {
		return "", err
	}

	return msg.FormatExecResponseMessage(&msg.ExecResponseMessageParams{
		ExitCode: exitCode,
		Output:   summary,
		Summary:  tail,
	})
}
