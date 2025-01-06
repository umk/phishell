package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type SummaryMessageParams struct {
	Summary string
}

//go:embed templates/summary_message.tmpl
var summaryMessage string

var summaryMessageTmpl = template.Must(template.New("summary_message").Parse(summaryMessage))

func FormatSummaryMessage(params *SummaryMessageParams) (string, error) {
	var sb strings.Builder
	if err := summaryMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
