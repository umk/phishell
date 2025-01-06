package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type ExecSummaryMessageParams struct {
	CommandLine string
	ExitCode    int
	Output      string
	Tail        bool // Indicates whether not full output is provided, but just a tail
}

//go:embed templates/exec_summary_message.tmpl
var execSummaryMessage string

var execSummaryMessageTmpl = template.Must(template.New("exec_summary_message").Parse(execSummaryMessage))

func FormatExecSummaryMessage(params *ExecSummaryMessageParams) (string, error) {
	var sb strings.Builder
	if err := execSummaryMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
