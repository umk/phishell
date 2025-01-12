package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type UserMessageParams struct {
	Request string

	Knowledge string
	Context   []string
}

//go:embed templates/user_message.tmpl
var userMessage string

var userMessageTmpl = template.Must(template.New("user_message").Parse(userMessage))

func FormatUserMessage(params *UserMessageParams) (string, error) {
	var sb strings.Builder
	if err := userMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
