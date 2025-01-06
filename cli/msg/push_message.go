package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type PushMessageParams struct {
	CommandLine string
	Output      string
}

//go:embed templates/push_message.tmpl
var pushMessage string

var pushMessageTmpl = template.Must(template.New("push_message").Parse(pushMessage))

func FormatPushMessage(params *PushMessageParams) (string, error) {
	var sb strings.Builder
	if err := pushMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
