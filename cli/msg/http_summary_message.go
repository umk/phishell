package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type HttpSummaryMessageParams struct {
	Url     string
	Status  string
	Headers map[string][]string
	Body    string
}

//go:embed templates/http_summary_message.tmpl
var httpSummaryMessage string

var httpSummaryMessageTmpl = template.Must(template.New("http_summary_message").Parse(httpSummaryMessage))

func FormatHttpSummaryMessage(params *HttpSummaryMessageParams) (string, error) {
	var sb strings.Builder
	if err := httpSummaryMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
