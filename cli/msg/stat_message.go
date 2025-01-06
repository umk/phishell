package msg

import (
	_ "embed"
	"strings"
	"text/template"
)

type StatMessageParams struct {
	IsDirectory bool
	Size        int64
}

//go:embed templates/stat_message.tmpl
var statMessage string

var statMessageTmpl = template.Must(template.New("stat_message").Parse(statMessage))

func FormatStatMessage(params *StatMessageParams) (string, error) {
	var sb strings.Builder
	if err := statMessageTmpl.Execute(&sb, params); err != nil {
		return "", err
	}

	return strings.TrimSpace(sb.String()), nil
}
