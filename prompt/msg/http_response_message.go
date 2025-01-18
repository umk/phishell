package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type HttpResponseMessageParams struct {
	Status  string
	Headers map[string][]string
	Body    string
	Summary bool
}

//go:embed templates/http_response_message.tmpl
var httpResponseMessage string

var httpResponseMessageTmpl = template.Must(template.New("http_response_message").Parse(httpResponseMessage))

func FormatHttpResponseMessage(params *HttpResponseMessageParams) (string, error) {
	var sb strings.Builder
	if err := httpResponseMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
