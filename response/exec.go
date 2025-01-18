package response

import (
	"context"
	"fmt"

	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/prompt"
	"github.com/umk/phishell/prompt/msg"
	"github.com/umk/phishell/util/execx"
)

type ExecOutputParams struct {
	CommandLine string
	ExitCode    int
	Output      execx.ProcessOutput
}

func GetExecOutput(ctx context.Context, cr *bootstrap.ClientRef, params *ExecOutputParams) (string, error) {
	outputStr, tail, err := params.Output.Get()
	if err != nil {
		return "", fmt.Errorf("invalid output: %w", err)
	}

	if params.ExitCode == 0 {
		if s, ok := getJSONObjectOrArray(outputStr); ok {
			outputStr = s
		}
	}

	if !mustSummarizeResp(cr, outputStr) {
		return msg.FormatExecResponseMessage(&msg.ExecResponseMessageParams{
			ExitCode: params.ExitCode,
			Output:   outputStr,
			Tail:     tail,
			Summary:  false,
		})
	}

	summaryCompl, err := prompt.PromptExecSummary(ctx, &prompt.ExecSummaryPromptParams{
		Client:      cr,
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
		Tail:     false,
		Summary:  true,
	})
}
