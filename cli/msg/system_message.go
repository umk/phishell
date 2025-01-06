package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type SystemMessageParams struct {
	Prompt string
	OS     string
}

//go:embed templates/system_message.tmpl
var systemMessage string

var systemMessageTmpl = template.Must(template.New("system_message").Parse(systemMessage))

func FormatSystemMessage(params *SystemMessageParams) (string, error) {
	var sb strings.Builder
	if err := systemMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
