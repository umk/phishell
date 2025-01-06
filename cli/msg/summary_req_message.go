package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type SummaryReqMessageParams struct{}

//go:embed templates/summary_req_message.tmpl
var summaryReqMessage string

var summaryReqMessageTmpl = template.Must(template.New("summary_req_message").Parse(summaryReqMessage))

func FormatSummaryReqMessage(params *SummaryReqMessageParams) (string, error) {
	var sb strings.Builder
	if err := summaryReqMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
