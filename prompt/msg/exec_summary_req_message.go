package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type ExecSummaryReqMessageParams struct {
	CommandLine string
	ExitCode    int
	Output      string
	Tail        bool // Indicates whether not full output is provided, but just a tail
}

//go:embed templates/exec_summary_req_message.tmpl
var execSummaryReqMessage string

var execSummaryReqMessageTmpl = template.Must(template.New("exec_summary_req_message").Parse(execSummaryReqMessage))

func FormatExecSummaryReqMessage(params *ExecSummaryReqMessageParams) (string, error) {
	var sb strings.Builder
	if err := execSummaryReqMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
