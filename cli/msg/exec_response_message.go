package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type ExecResponseMessageParams struct {
	ExitCode int
	Output   string
	Tail     bool
	Summary  bool
}

//go:embed templates/exec_response_message.tmpl
var execResponseMessage string

var execResponseMessageTmpl = template.Must(template.New("exec_response_message").Parse(execResponseMessage))

func FormatExecResponseMessage(params *ExecResponseMessageParams) (string, error) {
	var sb strings.Builder
	if err := execResponseMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
